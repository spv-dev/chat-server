package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

const (
	address = "localhost:50061"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	client := desc.NewChatV1Client(conn)

	// Создаем новый чат на сервере
	chat, err := createChat(ctx, client)
	if err != nil {
		log.Fatalf("failed to create chat: %v", err)
	}

	log.Printf(fmt.Sprintf("%s: %s\n", color.GreenString("Chat created"), color.YellowString(chat.Info.GetTitle())))

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Подключаемся к чату от имени пользователя oleg
	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chat.GetId(), 100500, 5*time.Second)
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	// Подключаемся к чату от имени пользователя ivan
	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chat.GetId(), 5500, 7*time.Second)
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	wg.Wait()
}

func connectChat(ctx context.Context, client desc.ChatV1Client, chatID int64, userID int64, period time.Duration) error {
	stream, err := client.ConnectChat(ctx, &desc.ConnectChatRequest{
		ChatId: chatID,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Println("failed to receive message from stream: ", errRecv)
				return
			}
			if userID == message.Info.UserId {
				continue
			}

			log.Printf("[%s] - [from: %s]: %s\n",
				color.YellowString(message.GetCreatedAt().AsTime().Format(time.RFC3339)),
				color.BlueString(strconv.FormatInt(message.Info.GetUserId(), 10)),
				message.Info.GetBody(),
			)
		}
	}()

	for {
		// Ниже пример того, как можно считывать сообщения из консоли
		// в демонстрационных целях будем засылать в чат рандомный текст раз в 5 секунд
		//scanner := bufio.NewScanner(os.Stdin)
		//var lines strings.Builder
		//
		//for {
		//	scanner.Scan()
		//	line := scanner.Text()
		//	if len(line) == 0 {
		//		break
		//	}
		//
		//	lines.WriteString(line)
		//	lines.WriteString("\n")
		//}
		//
		//err = scanner.Err()
		//if err != nil {
		//	log.Println("failed to scan message: ", err)
		//}

		time.Sleep(period)

		text := gofakeit.Word()

		_, err = client.SendMessage(ctx, &desc.SendMessageRequest{
			Info: &desc.MessageInfo{
				ChatId: chatID,
				UserId: userID,
				Body:   text,
			},
		})
		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}
}

func createChat(ctx context.Context, client desc.ChatV1Client) (*desc.Chat, error) {
	title := gofakeit.JobTitle()
	res, err := client.CreateChat(ctx, &desc.CreateChatRequest{
		Info: &desc.ChatInfo{
			Title: title,
		},
	})
	if err != nil {
		return &desc.Chat{}, err
	}

	return res.Chat, nil
}
