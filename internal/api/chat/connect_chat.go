package chat

import (
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ConnectChat соединение с сервером и получение сообщений
func (s *Server) ConnectChat(req *desc.ConnectChatRequest, stream desc.ChatV1_ConnectChatServer) error {
	s.mxChannel.RLock()
	chatChan, ok := s.channels[req.GetChatId()]
	s.mxChannel.RUnlock()

	if !ok {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	s.mxChat.Lock()
	if _, okChat := s.chats[req.GetChatId()]; !okChat {
		s.chats[req.GetChatId()] = &Chat{
			streams: make(map[int64]desc.ChatV1_ConnectChatServer),
		}
	}
	s.mxChat.Unlock()

	s.chats[req.GetChatId()].m.Lock()
	s.chats[req.GetChatId()].streams[req.GetUserId()] = stream
	s.chats[req.GetChatId()].m.Unlock()

	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}

			for _, st := range s.chats[req.GetChatId()].streams {
				if msg.Info.UserId != req.GetUserId() {
					if err := st.Send(msg); err != nil {
						return err
					}
				}
			}
		case <-stream.Context().Done():
			s.chats[req.GetChatId()].m.Lock()
			delete(s.chats[req.GetChatId()].streams, req.GetUserId())
			s.chats[req.GetChatId()].m.Unlock()
			return nil
		}
	}
}
