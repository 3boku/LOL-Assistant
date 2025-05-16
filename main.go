package main

import (
	"LOL-Assistant/disocrd"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// main 함수는 프로그램의 진입점이며 디스코드 봇을 초기화하고 실행합니다.
// .env 파일에서 환경 변수를 로드하고, Gemini 클라이언트를 초기화한 후
// 디스코드 봇을 시작합니다. 봇이 실행된 후에는 CTRL-C와 같은 인터럽트 신호를
// 기다려 프로그램을 종료합니다.
func main() {
	// .env 파일에서 환경 변수 로드
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}

	// Gemini 클라이언트 초기화
	disocrd.Initialize()

	// 디스코드 봇 토큰 설정
	token := fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN"))
	// 디스코드 세션 생성
	dg, err := discordgo.New(token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// 메시지 핸들러 등록
	dg.AddHandler(disocrd.Message)
	// 길드 메시지 인텐트 설정
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// 디스코드 연결 시작
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// 봇 실행 중 메시지 출력
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// 종료 신호 대기를 위한 채널 생성
	sc := make(chan os.Signal, 1)
	// 종료 신호 감지 설정
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// 신호를 수신할 때까지 블록
	<-sc

	// 디스코드 연결 종료
	defer func() {
		err = dg.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
