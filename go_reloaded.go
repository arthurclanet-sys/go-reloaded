package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const fichierEntree = "hello.txt"
const fichierSortie = "output.txt"

// commenceParVoyelleOuH renvoie true si le mot commence par a, e, i, o, u ou h.
func commenceParVoyelleOuH(mot string) bool {
	if mot == "" {
		return false
	}
	premiere, _ := utf8.DecodeRuneInString(mot)
	return strings.ContainsRune("aeiouh", unicode.ToLower(premiere))
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

// appliquerMarqueurs gère (hex), (bin), (up), (low) et (cap).
// Les marqueurs modifient le mot précédent et ne sont pas gardés dans le résultat.
func appliquerMarqueurs(mots []string) []string {
	var result []string

	for _, mot := range mots {
		dernier := len(result) - 1

		switch mot {
		case "(hex)", "(bin)", "(up)", "(low)", "(cap)":
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

		default:
			result = append(result, mot)
		}
	}
	return result
}

// collerApostrophes colle les ' autour des mots qu'ils entourent.
// " ' awesome ' " devient "'awesome'".
func collerApostrophes(mots []string) []string {
	var result []string
	ouvert := false          // true si on est entre deux '
	collerAuSuivant := false // true si le prochain mot doit recevoir un ' devant

	for _, mot := range mots {
		if mot == "'" {
			if !ouvert {
				// Apostrophe ouvrante : elle sera collée au mot suivant
				collerAuSuivant = true
				ouvert = true
			} else {
				// Apostrophe fermante : on la colle au mot précédent
				if len(result) > 0 {
					result[len(result)-1] += "'"
				} else {
					result = append(result, "'")
				}
				ouvert = false
			}
			continue
		}

		if collerAuSuivant {
			mot = "'" + mot
			collerAuSuivant = false
		}
		result = append(result, mot)
	}

	// Cas d'une apostrophe ouvrante seule tout à la fin du texte
	if collerAuSuivant {
		result = append(result, "'")
	}
	return result
}

// collerPonctuation colle . , ! ? : ; au mot précédent.
// "there ,and" devient "there, and" et "BAMM !!" devient "BAMM!!".
// Les groupes comme ... ou !? restent ensemble.
func collerPonctuation(mots []string) []string {
	var result []string

	for _, mot := range mots {
		// On compte combien de signes de ponctuation il y a au début du mot
		i := 0
		for i < len(mot) && strings.ContainsRune(".,!?:;", rune(mot[i])) {
			i++
		}
		ponctuation := mot[:i]
		reste := mot[i:]

		// La ponctuation est collée au mot précédent
		if ponctuation != "" {
			if len(result) > 0 {
				result[len(result)-1] += ponctuation
			} else {
				result = append(result, ponctuation)
			}
		}

		// Le reste du mot (s'il y en a) devient un mot séparé
		if reste != "" {
			result = append(result, reste)
		}
	}
	return result
}

// transformerAEnAn remplace "a" par "an" si le mot suivant commence par une voyelle ou un h.
func transformerAEnAn(mots []string) []string {
	for i := 0; i < len(mots)-1; i++ {
		if commenceParVoyelleOuH(mots[i+1]) {
			if mots[i] == "a" {
				mots[i] = "an"
			} else if mots[i] == "A" {
				mots[i] = "An"
			}
		}
	}
	return mots
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

	// Étape 3 : appliquer toutes les règles, dans l'ordre
	result := appliquerMarqueurs(words)
	result = collerApostrophes(result)
	result = collerPonctuation(result)
	result = transformerAEnAn(result)

	// Étape 4 : recoller les mots et écrire le résultat
	final := strings.Join(result, " ") + "\n"

	err = os.WriteFile(fichierSortie, []byte(final), 0644)
	if err != nil {
		fmt.Println("Erreur : impossible d'écrire", fichierSortie)
		os.Exit(1)
	}

	fmt.Print("Résultat : ", final)
	fmt.Println("Fichier", fichierSortie, "créé avec succès !")
}
