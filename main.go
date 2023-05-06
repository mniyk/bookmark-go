package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	bookmarks = []Bookmark{}
	file_name = "bookmark.json"
)

type Bookmark struct {
	Title string
	URL   string
}

func read_bookmarks() {
	// jsonファイルの読み込み
	byte_data, err := os.ReadFile(file_name)
	if err != nil {
		// jsonファイルの作成
		file, err := os.Create(file_name)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	// jsonデータの構造体インスタンス化
	_ = json.Unmarshal(byte_data, &bookmarks)
}

func add_bookmark(bookmark Bookmark) {
	// read_bookmarks()

	exist := false
	for _, b := range bookmarks {
		if b.URL == bookmark.URL {
			exist = true
		}
	}

	if !exist {
		bookmarks = append(bookmarks, bookmark)
	}

	write_json_data()
}

func delete_bookmark(bookmark Bookmark) {
	new_bookmarks := []Bookmark{}

	for _, b := range bookmarks {
		if b.URL != bookmark.URL {
			new_bookmarks = append(new_bookmarks, b)
		}
	}

	bookmarks = new_bookmarks

	write_json_data()
}

func write_json_data() {
	// jsonファイルの取得
	file, _ := os.Create(file_name)
	defer file.Close()

	// jsonファイルへの書き込み
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(bookmarks); err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	// Form
	titleInput := tview.NewInputField().SetLabel("Title").SetText("")
	urlInput := tview.NewInputField().SetLabel("URL").SetText("")
	form := tview.NewForm().
		AddFormItem(titleInput).
		AddFormItem(urlInput).
		AddButton("Save", nil).
		AddButton("Delete", nil)
	form.SetBorder(true).SetTitle(" Bookmark Form ").SetTitleAlign(tview.AlignLeft)
	form.SetFieldBackgroundColor(tcell.ColorGreen).SetButtonBackgroundColor(tcell.ColorGreen).SetLabelColor(tcell.ColorWhite)

	// Table
	list := tview.NewTable().SetFixed(1, 1)
	read_bookmarks()
	h_title := tview.NewTableCell("Title").SetSelectable(false)
	h_url := tview.NewTableCell("URL").SetSelectable(false)
	list.SetCell(0, 0, h_title)
	list.SetCell(0, 1, h_url)
	for i, bookmark := range bookmarks {
		title := tview.NewTableCell(bookmark.Title)
		url := tview.NewTableCell(bookmark.URL)
		list.SetCell(i+1, 0, title)
		list.SetCell(i+1, 1, url)
	}
	list.SetBorder(true).SetTitle(" Bookmark List ").SetTitleAlign(tview.AlignLeft)
	list.SetSelectable(true, false).SetSeparator(' ')

	// Frame
	frame := tview.NewFrame(tview.NewBox()).
		AddText("Move Bookmark Form: Ctrl + f", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("Move Bookmark List: Ctrl + l", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("URL Open: Ctrl + o (selected bookmark)", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("Edit Bookmark: Ctrl + e (selected bookmark)", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("Quit: Ctrl + c", true, tview.AlignLeft, tcell.ColorWhite)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(frame, 8, 1, false).
		AddItem(form, 9, 1, false).
		AddItem(list, 0, 1, true)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlO {
			row, _ := list.GetSelection()
			url := list.GetCell(row, 1).Text
			if err := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url).Start(); err != nil {
				log.Fatalln(err)
			}
		}

		return event
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlE {
			row, _ := list.GetSelection()
			title := list.GetCell(row, 0).Text
			url := list.GetCell(row, 1).Text
			titleInput.SetText(title)
			urlInput.SetText(url)
			app.ForceDraw()
		} else if event.Key() == tcell.KeyCtrlL {
			app.SetFocus(list)
		} else if event.Key() == tcell.KeyCtrlF {
			app.SetFocus(form)
		}

		return event
	})

	form.GetButton(0).SetSelectedFunc(func() {
		list.Clear()

		list.SetCell(0, 0, h_title)
		list.SetCell(0, 1, h_url)

		bookmark := Bookmark{
			Title: titleInput.GetText(),
			URL:   urlInput.GetText(),
		}

		add_bookmark(bookmark)

		for i, b := range bookmarks {
			title := tview.NewTableCell(b.Title)
			url := tview.NewTableCell(b.URL)
			list.SetCell(i+1, 0, title)
			list.SetCell(i+1, 1, url)
		}

		titleInput.SetText("")
		urlInput.SetText("")

		app.ForceDraw()
	})

	form.GetButton(1).SetSelectedFunc(func() {
		list.Clear()

		list.SetCell(0, 0, h_title)
		list.SetCell(0, 1, h_url)

		bookmark := Bookmark{
			Title: titleInput.GetText(),
			URL:   urlInput.GetText(),
		}

		delete_bookmark(bookmark)

		for i, b := range bookmarks {
			title := tview.NewTableCell(b.Title)
			url := tview.NewTableCell(b.URL)
			list.SetCell(i+1, 0, title)
			list.SetCell(i+1, 1, url)
		}

		titleInput.SetText("")
		urlInput.SetText("")

		app.ForceDraw()
	})

	pages.AddPage("main", layout, true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
