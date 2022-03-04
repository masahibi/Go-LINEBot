package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	//"strconv"
	"strings"
)

// タスク帳の項目
type Task struct { // TaskBookのタスクを表す構造体
	Category string
	Date     string
}

// 家計簿の処理を行う型
type TaskBook struct { // TaskBookを表す構造体
	file  string
	tasks []*Task
}

// 新しいTaskBookを作成する
func NewTaskBook(file string) *TaskBook {
	// TaskBook構造体を作成する
	taskBook := &TaskBook{
		file: file,
	}

	taskBook.readItems()

	// TaskBookのポインタを返す
	return taskBook
}

func (taskBook *TaskBook) readItems() {
	f, err := os.Open(taskBook.file) // ファイルを読み込みで開く
	//if err != nil {  // エラーが発生したら
	//	fmt.Fprintln(os.Stderr, "エラー：", err) // 標準エラー出力に出力
	//	os.Exit(1)                           // 終了コードを指定してプログラムを終了
	//}
	if err != nil { // エラーが発生した場合
		log.Fatal(err) // エラー内容を出力して終了
	}

	defer f.Close() // ファイルを閉じる (遅延実行)

	s := bufio.NewScanner(f) // ファイルから読み込む
	for s.Scan() {           // 1行ずつ読み込む
		ss := strings.Split(s.Text(), ",") // ',' 区切りで分割してスライスにする
		if len(ss) != 2 {
			fmt.Fprintln(os.Stderr, "ファイル形式が不正です") // 標準エラー出力に出力
			os.Exit(1)                             // 終了コードを指定してプログラムを終了
		}

		//date, err := strconv.Atoi(ss[1]) // 文字列をint型に変換  <----------------------------------------------- 後ほど
		date := ss[1]

		//if err != nil {
		//	fmt.Fprintln(os.Stderr, "エラー：", err) // 標準エラー出力に出力
		//	os.Exit(1)                           // 終了コードを指定してプログラムを終了
		//}
		if err != nil { // エラーが発生した場合
			log.Fatal(err) // エラー内容を出力して終了
		}

		task := &Task{ // Task構造体を作成
			Category: ss[0],
			Date:     date,
		}

		taskBook.AddTask(task) // タスクを追加する
	}

	//if err := s.Err(); err != nil { // まとめてエラー処理
	//	fmt.Fprintln(os.Stderr, "エラー：", err) // 標準エラー出力に出力
	//	os.Exit(1)                           // 終了コードを指定してプログラムを終了
	//}
	if err := s.Err(); err != nil { // エラーが発生した場合
		log.Fatal(err) // エラー内容を出力して終了
	}

}

// 新しいTaskを追加する
func (taskBook *TaskBook) AddTask(task *Task) {
	taskBook.tasks = append(taskBook.tasks, task) // 追加する

	f, err := os.Create(taskBook.file) // ファイルを書込みで開く
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "エラー：", err) // 標準エラー出力に出力
	//	os.Exit(1)                           // 終了コードを指定してプログラムを終了
	//}
	if err != nil { // エラーが発生した場合
		log.Fatal(err) // エラー内容を出力して終了
	}

	for _, task := range taskBook.tasks { // タスクを1つずつ取り出す
		_, err := fmt.Fprintf(f, "%s,%s\n", task.Category, task.Date) // ファイルに書き込む
		//if err != nil {
		//	fmt.Fprintln(os.Stderr, "エラー：", err) // 標準エラー出力に出力
		//	os.Exit(1)                           // 終了コードを指定してプログラムを終了
		//}
		if err != nil { // エラーが発生した場合
			log.Fatal(err) // エラー内容を出力して終了
		}
	}

	//if err := f.Close(); err != nil {
	//	fmt.Fprintln(os.Stderr, "エラー：", err) // 標準エラー出力に出力
	//	os.Exit(1)                           // 終了コードを指定してプログラムを終了
	//}
	if err := f.Close(); err != nil {
		log.Fatal(err) // エラー内容を出力して終了
	}
}
