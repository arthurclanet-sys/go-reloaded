package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	// Étape 1 : lire le fichier
	content, err := os.ReadFile("hello.txt")
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	// Recopier le contenu dans un autre fichier
	err = os.WriteFile("output.txt", content, 0644)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("Fichier copié avec succès !")

	// Étape 2 : découper le texte en mots
	text := string(content)

	words := strings.Fields(text)

	// Étape 3 : convertir les nombres hexadécimaux et binaires
	for i := 0; i < len(words); i++ {

		// Conversion HEX -> DECIMAL
		if words[i] == "(hex)" && i > 0 {
			n, err := strconv.ParseInt(words[i-1], 16, 64)
			if err == nil {
				words[i-1] = strconv.FormatInt(n, 10)
			}
		}

		// Conversion BINAIRE -> DECIMAL
		if words[i] == "(bin)" && i > 0 {
			n, err := strconv.ParseInt(words[i-1], 2, 64)
			if err == nil {
				words[i-1] = strconv.FormatInt(n, 10)
			}
		}
	}

	fmt.Println("\nListe des mots après conversion :")

	for i := 0; i < len(words); i++ {
		fmt.Println(words[i])
	}
}
