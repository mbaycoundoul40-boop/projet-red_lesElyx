package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

type Equipment struct{ Head, Chest, Feet string }
type Character struct {
	Name, Class                                                string
	Level, MaxHP, CurrentHP, Gold, InventoryLimit, InvUpgrades int
	Inventory, Skills                                          []string
	Equipment                                                  Equipment
}
type Monster struct {
	Name                     string
	MaxHP, CurrentHP, Attack int
}

func initCharacter(name, class string, level, maxHP, cur int, inv []string) Character {
	return Character{Name: name, Class: class, Level: level, MaxHP: maxHP, CurrentHP: cur, Inventory: inv, InventoryLimit: 10, Gold: 100, Skills: []string{"Onde"}}
}
func displayInfo(c Character) {
	fmt.Println("=== Fiche de", c.Name, "=== Classe :", c.Class, " Niveau :", c.Level)
	fmt.Println("PV :", c.CurrentHP, "/", c.MaxHP, " Or :", c.Gold, " Pouvoirs :", c.Skills)
	fmt.Println("Équipement :", c.Equipment, " Inventaire :", c.Inventory)
}
func canAddItem(c *Character) bool { return len(c.Inventory) < c.InventoryLimit }
func addInventory(c *Character, item string) bool {
	if !canAddItem(c) {
		fmt.Println("Inventaire plein ! Impossible d'ajouter :", item)
		return false
	}
	c.Inventory = append(c.Inventory, item)
	return true
}
func removeInventory(c *Character, item string) bool {
	for i, it := range c.Inventory {
		if it == item {
			c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			return true
		}
	}
	return false
}
func countItem(c *Character, item string) int {
	n := 0
	for _, it := range c.Inventory {
		if it == item {
			n++
		}
	}
	return n
}
func takePot(c *Character) {
	if !removeInventory(c, "Fiole d'Énergie") {
		fmt.Println("Vous n'avez pas de Fiole d'Énergie !")
		return
	}
	c.CurrentHP += 50
	if c.CurrentHP > c.MaxHP {
		c.CurrentHP = c.MaxHP
	}
	fmt.Println("Fiole utilisée. PV :", c.CurrentHP, "/", c.MaxHP)
}
func isDead(c *Character) bool {
	if c.CurrentHP <= 0 {
		c.CurrentHP = c.MaxHP / 2
		fmt.Println(c.Name, "vaincu(e) ! Réveil avec", c.CurrentHP, "/", c.MaxHP, "PV.")
		return true
	}
	return false
}
func poisonPot(c *Character) {
	if !removeInventory(c, "Fiole Corrompue") {
		fmt.Println("Vous n'avez pas de Fiole Corrompue !")
		return
	}
	for i := 0; i < 3; i++ {
		c.CurrentHP -= 10
		if c.CurrentHP < 0 {
			c.CurrentHP = 0
		}
		fmt.Println("PV :", c.CurrentHP, "/", c.MaxHP)
		time.Sleep(time.Second)
	}
	isDead(c)
}
func spellBook(c *Character) {
	for _, s := range c.Skills {
		if s == "Déflagration" {
			fmt.Println("Déjà appris.")
			return
		}
	}
	c.Skills = append(c.Skills, "Déflagration")
	fmt.Println("Appris : Déflagration !")
}
func accessInventory(c *Character) {
	for {
		fmt.Println("=== Inventaire ===")
		for i, item := range c.Inventory {
			fmt.Println(i+1, "-", item)
		}
		fmt.Println("0 - Retour")
		var ch int
		fmt.Scan(&ch)
		if ch == 0 {
			return
		}
		if ch < 1 || ch > len(c.Inventory) {
			fmt.Println("Choix invalide.")
			continue
		}
		switch item := c.Inventory[ch-1]; item {
		case "Fiole d'Énergie":
			takePot(c)
		case "Fiole Corrompue":
			poisonPot(c)
		case "Grimoire : Déflagration":
			spellBook(c)
			removeInventory(c, item)
		case "Heaume du Zénith", "Plastron du Zénith", "Bottes du Zénith":
			equipItem(c, item)
		default:
			fmt.Println("Inutilisable ici.")
		}
	}
}
func formatName(s string) string { s = strings.ToLower(s); return strings.ToUpper(s[:1]) + s[1:] }
func onlyLetters(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
func characterCreation() Character {
	var name string
	for {
		fmt.Println("Nom (lettres uniquement) ?")
		fmt.Scan(&name)
		if onlyLetters(name) {
			break
		}
		fmt.Println("Invalide.")
	}
	classes := map[int][2]interface{}{1: {"Humaine", 100}, 2: {"Elfe", 80}, 3: {"Naine", 120}}
	var ch int
	for {
		fmt.Println("Classe : 1-Humaine(100) 2-Elfe(80) 3-Naine(120)")
		fmt.Scan(&ch)
		if _, ok := classes[ch]; ok {
			break
		}
		fmt.Println("Invalide.")
	}
	class, maxHP := classes[ch][0].(string), classes[ch][1].(int)
	c := initCharacter(formatName(name), class, 1, maxHP, maxHP/2, []string{})
	fmt.Println("Bienvenue,", c.Name, "! L'ascension du Zénith commence.")
	return c
}
func merchant(c *Character) {
	items := map[int]struct {
		nom  string
		cout int
	}{1: {"Fiole d'Énergie", 3}, 2: {"Fiole Corrompue", 6}, 3: {"Grimoire : Déflagration", 25}, 4: {"Écaille de Golem", 4}, 5: {"Griffe d'Ombre", 7}, 6: {"Cuir de Bête", 3}, 7: {"Plume Astrale", 1}}
	for {
		fmt.Println("=== Marchande === Or :", c.Gold)
		for i := 1; i <= 7; i++ {
			fmt.Println(i, "-", items[i].nom, "(", items[i].cout, "or)")
		}
		fmt.Println("8 - Augmentation inventaire (30 or)  0 - Retour")
		var ch int
		fmt.Scan(&ch)
		if ch == 0 {
			return
		}
		if ch == 8 {
			if c.Gold < 30 {
				fmt.Println("Pas assez d'or.")
			} else if c.InvUpgrades >= 3 {
				fmt.Println("Limite atteinte (3).")
			} else {
				c.Gold -= 30
				upgradeInventorySlot(c)
			}
			continue
		}
		art, ok := items[ch]
		if !ok {
			fmt.Println("Invalide.")
			continue
		}
		if c.Gold < art.cout {
			fmt.Println("Pas assez d'or.")
		} else if !canAddItem(c) {
			fmt.Println("Inventaire plein.")
		} else {
			c.Gold -= art.cout
			addInventory(c, art.nom)
			fmt.Println("Acheté :", art.nom)
		}
	}
}
func blacksmith(c *Character) {
	type r struct {
		out  string
		ingr map[string]int
	}
	recipes := map[int]r{1: {"Heaume du Zénith", map[string]int{"Plume Astrale": 1, "Cuir de Bête": 1}}, 2: {"Plastron du Zénith", map[string]int{"Écaille de Golem": 2, "Griffe d'Ombre": 1}}, 3: {"Bottes du Zénith", map[string]int{"Écaille de Golem": 1, "Cuir de Bête": 1}}}
	for {
		fmt.Println("=== Forge === Or :", c.Gold)
		fmt.Println("1-Heaume(5or+1Plume+1Cuir) 2-Plastron(5or+2Écaille+1Griffe) 3-Bottes(5or+1Écaille+1Cuir) 0-Retour")
		var ch int
		fmt.Scan(&ch)
		if ch == 0 {
			return
		}
		rc, ok := recipes[ch]
		if !ok {
			fmt.Println("Invalide.")
			continue
		}
		if c.Gold < 5 || !canAddItem(c) {
			fmt.Println("Or ou place insuffisants.")
			continue
		}
		miss := false
		for it, q := range rc.ingr {
			if countItem(c, it) < q {
				miss = true
			}
		}
		if miss {
			fmt.Println("Ressources manquantes.")
			continue
		}
		for it, q := range rc.ingr {
			for i := 0; i < q; i++ {
				removeInventory(c, it)
			}
		}
		c.Gold -= 5
		addInventory(c, rc.out)
		fmt.Println("Fabriqué :", rc.out)
	}
}
func equipBonus(item string) int {
	switch item {
	case "Heaume du Zénith":
		return 10
	case "Plastron du Zénith":
		return 25
	case "Bottes du Zénith":
		return 15
	}
	return 0
}
func equipSlot(c *Character, item string) *string {
	switch item {
	case "Heaume du Zénith":
		return &c.Equipment.Head
	case "Plastron du Zénith":
		return &c.Equipment.Chest
	case "Bottes du Zénith":
		return &c.Equipment.Feet
	}
	return nil
}
func equipItem(c *Character, item string) {
	slot := equipSlot(c, item)
	if slot == nil {
		return
	}
	if *slot != "" {
		addInventory(c, *slot)
		c.MaxHP -= equipBonus(*slot)
	}
	removeInventory(c, item)
	*slot = item
	c.MaxHP += equipBonus(item)
	if c.CurrentHP > c.MaxHP {
		c.CurrentHP = c.MaxHP
	}
	fmt.Println("Équipé :", item, "(PV max :", c.MaxHP, ")")
}
func upgradeInventorySlot(c *Character) {
	c.InventoryLimit += 10
	c.InvUpgrades++
	fmt.Println("Nouvelle limite :", c.InventoryLimit)
}
func initGuardian() Monster { return Monster{Name: "Gardien du Seuil", MaxHP: 40, CurrentHP: 40, Attack: 5} }
func guardianPattern(m *Monster, c *Character, turn int) {
	dmg := m.Attack
	if turn%3 == 0 {
		dmg *= 2
	}
	c.CurrentHP -= dmg
	if c.CurrentHP < 0 {
		c.CurrentHP = 0
	}
	fmt.Println(m.Name, "inflige à", c.Name, dmg, "dégâts. PV :", c.CurrentHP, "/", c.MaxHP)
}
func characterTurn(c *Character, m *Monster) {
	hasSpell := false
	for _, s := range c.Skills {
		if s == "Déflagration" {
			hasSpell = true
		}
	}
	opts := "1-Attaquer 2-Inventaire"
	if hasSpell {
		opts += " 3-Déflagration"
	}
	fmt.Println("=== Tour de", c.Name, "===", opts)
	var ch int
	fmt.Scan(&ch)
	switch ch {
	case 1:
		dmg := 5
		m.CurrentHP -= dmg
		if m.CurrentHP < 0 {
			m.CurrentHP = 0
		}
		fmt.Println(c.Name, "inflige", dmg, "dégâts à", m.Name, "PV :", m.CurrentHP, "/", m.MaxHP)
	case 2:
		usable := []string{}
		for _, item := range c.Inventory {
			if item == "Fiole d'Énergie" || item == "Fiole Corrompue" {
				usable = append(usable, item)
			}
		}
		if len(usable) == 0 {
			fmt.Println("Aucun objet utilisable.")
			return
		}
		for i, item := range usable {
			fmt.Println(i+1, "-", item)
		}
		var ic int
		fmt.Scan(&ic)
		if ic < 1 || ic > len(usable) {
			fmt.Println("Invalide.")
			return
		}
		if usable[ic-1] == "Fiole d'Énergie" {
			takePot(c)
		} else {
			poisonPot(c)
		}
	case 3:
		if !hasSpell {
			fmt.Println("Tour passé.")
			return
		}
		dmg := 15
		m.CurrentHP -= dmg
		if m.CurrentHP < 0 {
			m.CurrentHP = 0
		}
		fmt.Println(c.Name, "lance Déflagration :", dmg, "dégâts à", m.Name, "PV :", m.CurrentHP, "/", m.MaxHP)
	default:
		fmt.Println("Tour passé.")
	}
}
func trainingFight(c *Character) {
	m := initGuardian()
	turn := 1
	fmt.Println("Le Gardien du Seuil bloque votre passage !")
	for {
		fmt.Println("----- Tour", turn, "-----")
		characterTurn(c, &m)
		if m.CurrentHP <= 0 {
			c.Level++
			c.MaxHP += 10
			c.CurrentHP = c.MaxHP
			fmt.Println(m.Name, "vaincu ! Zénith niveau", c.Level, "atteint. PV max :", c.MaxHP)
			break
		}
		guardianPattern(&m, c, turn)
		if c.CurrentHP <= 0 {
			isDead(c)
			break
		}
		turn++
	}
	fmt.Println("Retour au menu.")
}
func menu(c *Character) {
	for {
		fmt.Println("\n=== Menu === 1-Infos 2-Inventaire 3-Marchande 4-Forge 5-Défi 6-Quitter")
		var ch int
		fmt.Scan(&ch)
		switch ch {
		case 1:
			displayInfo(*c)
		case 2:
			accessInventory(c)
		case 3:
			merchant(c)
		case 4:
			blacksmith(c)
		case 5:
			trainingFight(c)
		case 6:
			fmt.Println("À bientôt,", c.Name, "!")
			return
		default:
			fmt.Println("Invalide.")
		}
	}
}
func main() {
	fmt.Println("=== L'Ascension du Zénith ===")
	c1 := characterCreation()
	menu(&c1)
}
