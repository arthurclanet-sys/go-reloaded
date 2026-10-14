package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

const fichierEntree = "hello.txt"
const fichierSortie = "output.txt"

// commenceParVoyelle renvoie true si le mot commence par une voyelle (accents compris).
func commenceParVoyelle(mot string) bool {
	if mot == "" {
		return false
	}
	premiere, _ := utf8.DecodeRuneInString(strings.ToLower(mot))
	return strings.ContainsRune("aeiouyàâäéèêëîïôöùûü", premiere)
}

// convertir transforme un nombre écrit dans une base (16 ou 2) en décimal.
// Si le mot n'est pas un nombre valide, il est renvoyé tel quel avec un avertissement.
func convertir(mot string, base int) string {
	n, err := strconv.ParseInt(mot, base, 64)
	if err != nil {
		fmt.Println("Attention :", mot, "n'est pas un nombre valide en base", base, "- mot laissé tel quel.")
		return mot
	}
	return strconv.FormatInt(n, 10)
}

// capitaliser met la première lettre en majuscule et le reste en minuscules.
func capitaliser(mot string) string {
	if mot == "" {
		return mot
	}
	premiere, taille := utf8.DecodeRuneInString(mot)
	return strings.ToUpper(string(premiere)) + strings.ToLower(mot[taille:])
}

func main() {

	// Étape 1 : lire le fichier d'entrée
	content, err := os.ReadFile(fichierEntree)
	if err != nil {
		fmt.Println("Erreur : impossible de lire", fichierEntree)
		fmt.Println("Vérifie que le fichier existe bien dans le même dossier que main.go.")
		os.Exit(1)
	}

	// Étape 2 : découper le texte en mots
	words := strings.Fields(string(content))
	if len(words) == 0 {
		fmt.Println("Erreur :", fichierEntree, "est vide, il n'y a rien à transformer.")
		os.Exit(1)
	}

	// Étape 3 : appliquer les règles (hex), (bin), (up), (low), (cap)
	// Les marqueurs modifient le mot précédent et ne sont pas ajoutés au résultat.
	var result []string

	for _, mot := range words {
		dernier := len(result) - 1 // position du mot précédent dans result

		switch mot {
		case "(hex)", "(bin)", "(up)", "(low)", "(cap)":
			// Pas de mot avant le marqueur : on l'ignore
			if dernier < 0 {
				fmt.Println("Attention :", mot, "est placé au début du texte, il est ignoré.")
				continue
			}

			switch mot {
			case "(hex)":
				result[dernier] = convertir(result[dernier], 16)
			case "(bin)":
				result[dernier] = convertir(result[dernier], 2)
			case "(up)":
				result[dernier] = strings.ToUpper(result[dernier])
			case "(low)":
				result[dernier] = strings.ToLower(result[dernier])
			case "(cap)":
				result[dernier] = capitaliser(result[dernier])
			}

		// Mot normal : on l'ajoute tel quel
		default:
			result = append(result, mot)
		}
	}

	// Étape 4 : transformer "a" en "an" devant une voyelle
	for i := 0; i < len(result)-1; i++ {
		if commenceParVoyelle(result[i+1]) {
			if result[i] == "a" {
				result[i] = "an"
			} else if result[i] == "A" {
				result[i] = "An"
			}
		}
	}

	// Étape 5 : recoller les mots et écrire le résultat
	// (output.txt est créé automatiquement s'il n'existe pas)
	final := strings.Join(result, " ") + "\n"

	err = os.WriteFile(fichierSortie, []byte(final), 0644)
	if err != nil {
		fmt.Println("Erreur : impossible d'écrire", fichierSortie)
		os.Exit(1)
	}

	fmt.Print("Résultat : ", final)
}
