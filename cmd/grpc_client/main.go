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
	chatID, err := createChat(ctx, client)
	if err != nil {
		log.Fatalf("failed to create chat: %v", err)
	}

	log.Printf(fmt.Sprintf("%s: %s\n", color.GreenString("Chat created"), color.YellowString(strconv.FormatInt(chatID, 10))))

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Подключаемся к чату от имени пользователя oleg
	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chatID, 100500, 5*time.Second)
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	// Подключаемся к чату от имени пользователя ivan
	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chatID, 5500, 7*time.Second)
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

			log.Printf("[%s] - [from: %s]: %s\n",
				color.YellowString(strconv.FormatInt(message.GetChatId(), 10)),
				color.BlueString(strconv.FormatInt(message.GetUserId(), 10)),
				message.GetBody(),
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

func createChat(ctx context.Context, client desc.ChatV1Client) (int64, error) {
	res, err := client.CreateChat(ctx, &desc.CreateChatRequest{})
	if err != nil {
		return 0, err
	}

	return res.GetId(), nil
}
