package db

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {

	_, err := Db.Exec(`
        CREATE TABLE IF NOT EXISTS quiz (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
			created TIMESTAMP not null default current_timestamp,
            updated TIMESTAMP not null default current_timestamp,
            deleted TIMESTAMP,
            japanese TEXT UNIQUE NOT NULL,
            pronounce TEXT UNIQUE NOT NULL,
            english TEXT UNIQUE NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("error creating table quiz: %w", err)
	}

	// Migrate some data
	sentences := []struct {
		Japanese  string
		Pronounce string
		English   string
	}{
		{"こんにちは、世界！", "Konnichiwa, sekai!", "Hello, world!"},
		{"おはようございます。", "Ohayou gozaimasu.", "Good morning."},
		{"こんばんは。", "Konbanwa.", "Good evening."},
		{"ありがとう。", "Arigatou.", "Thank you."},
		{"ごめんなさい。", "Gomen nasai.", "I'm sorry."},
		{"はい、そうです。", "Hai, sou desu.", "Yes, that's right."},
		{"いいえ、違います。", "Iie, chigaimasu.", "No, that's wrong."},
		{"お願いします。", "Onegaishimasu.", "Please."},
		{"どういたしまして。", "Dou itashimashite.", "You're welcome."},
		{"いってきます。", "Ittekimasu.", "I'm leaving (said when leaving home)."},
		{"ただいま。", "Tadaima.", "I'm back (said when returning home)."},
		{"行ってらっしゃい。", "Itte rasshai.", "Take care (said when someone is leaving)."},
		{"お誕生日おめでとうございます！", "Otanjoubi omedetou gozaimasu!", "Happy birthday!"},
		{"おめでとうございます！", "Omedetou gozaimasu!", "Congratulations!"},
		{"おはよう、日本！", "Ohayou, Nihon!", "Good morning, Japan!"},
		{"こんな天気の日には外で遊びましょう。", "Konna tenki no hi ni wa soto de asobimashou.", "Let's play outside on such a sunny day."},
		{"もうすぐ春が来ます。", "Mou sugu haru ga kimasu.", "Spring is coming soon."},
		{"今日はとても寒いですね。", "Kyou wa totemo samui desu ne.", "It's very cold today, isn't it?"},
		{"明日は晴れるといいですね。", "Ashita wa hareru to ii desu ne.", "I hope it will be sunny tomorrow."},
		{"この本はとても面白いです。", "Kono hon wa totemo omoshiroi desu.", "This book is very interesting."},
	}

	stmt, err := Db.Prepare("INSERT INTO quiz (japanese, pronounce, english) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing db statements: %w", err)
	}
	defer stmt.Close()

	for _, s := range sentences {
		_, err := stmt.Exec(s.Japanese, s.Pronounce, s.English)
		if err != nil {
			return fmt.Errorf("error inserting data into db: %w", err)
		}
	}

	fmt.Println("Migration successful.")
	return nil
}
