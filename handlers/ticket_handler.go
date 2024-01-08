package handlers

import (
	"awesomeProject/appconst"
	"encoding/json"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"log"
	"net/http"
	"os"
)

type ticketHandler struct {
}

type TicketHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Confirm(w http.ResponseWriter, r *http.Request)
	Approve(w http.ResponseWriter, r *http.Request)
}

func NewTicketHandler() TicketHandler {
	return &ticketHandler{}
}

var targetQRCodeImage = fmt.Sprintf("%s/%s.png", appconst.QRCodeImageDir, appconst.TargetTicketToken)
var targetUserImage = fmt.Sprintf("%s/%d.jpeg", appconst.UserImageDir, appconst.TargetUserID)

func (*ticketHandler) Create(w http.ResponseWriter, r *http.Request) {
	// QRコードのデータを設定
	data := appconst.TargetTicketToken

	// QRコードを生成
	// 第２引数がQRコードのレベルを表す。L,M,Q,Hの４つがあり、それぞれ容量が異なる
	// 第３引数がQRコードのモードを表す。Auto,Binary,AlphaNumeric,Unicodeの４つがあり、それぞれ容量が異なる
	qrCode, err := qr.Encode(data, qr.M, qr.Auto)
	if err != nil {
		log.Fatal(err)
	}

	// QRコードのサイズを調整
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		log.Fatal(err)
	}

	// 画像ファイルを作成
	file, err := os.Create(fmt.Sprintf("./%s", targetQRCodeImage))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// PNGエンコーダーを使用してQRコードをファイルに保存
	if err := png.Encode(file, qrCode); err != nil {
		log.Fatal(err)
	}
}

func (*ticketHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ticketGet := struct {
		TicketURL string `json:"ticketURL"`
	}{
		TicketURL: fmt.Sprintf("http://localhost:8080/%s", targetQRCodeImage),
	}
	b, _ := json.Marshal(ticketGet)
	fmt.Fprintf(w, string(b))
}

func (*ticketHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ticketConfirm := struct {
		UserAge   int    `json:"userAge"`
		UserName  string `json:"userName"`
		UserImage string `json:"userImage"`
	}{
		UserAge:   36,
		UserName:  "曽根田侑也",
		UserImage: fmt.Sprintf("http://localhost:8080/%s", targetUserImage),
	}
	b, _ := json.Marshal(ticketConfirm)
	fmt.Fprintf(w, string(b))
}

func (*ticketHandler) Approve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ticketConfirm := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}{
		Count:  1,
		Status: "approved",
	}
	b, _ := json.Marshal(ticketConfirm)
	fmt.Fprintf(w, string(b))
}
