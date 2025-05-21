package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/cilium/tetragon/api/v1/tetragon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const tetragonAddress = "localhost:54321"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		tetragonAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("failed to connect to Tetragon: %v", err)
	}
	defer conn.Close()

	client := pb.NewFineGuidanceSensorsClient(conn)

	stream, err := client.GetEvents(context.Background(), &pb.GetEventsRequest{})
	if err != nil {
		log.Fatalf("GetEvents RPC failed: %v", err)
	}

	log.Println("Subscribed to Tetragon events — streaming now…")

	for {
		ev, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error receiving event: %v", err)
		}

		// Handle Kprobe events
		if k := ev.GetProcessKprobe(); k != nil {
			name := k.GetPolicyName()
			process := k.GetProcess()

			pid := process.GetPid().GetValue() >> 32
			comm := process.GetDocker()

			arg := k.GetArgs()[0].GetSockArg()
			saddr := arg.GetSaddr()
			daddr := arg.GetDaddr()

			ts := time.Unix(0, ev.GetTime().Seconds)

			fmt.Printf("[%s] KProbe %s: PID=%d COMM=%s saddr=%s daddr=%s\n",
				ts.Format(time.RFC3339),
				name,
				pid,
				comm,
				saddr,
				daddr,
			)
		}
	}
}
