package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/grpc"

	pb "github.com/dheeraj-sn/distributed-orchestrator/proto"
)

type job struct {
	ID     string
	Status string
	Result string
}

type model struct {
	jobs   []job
	cursor int
	err    error
	filter string
}

func initialModel() model {
	jobs, err := fetchJobs()
	return model{jobs: jobs, err: err, filter: ""}
}

func fetchJobs() ([]job, error) {
	addr := os.Getenv("SCHEDULER_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewOrchestratorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.ListJobs(ctx, &pb.ListJobsRequest{})
	if err != nil {
		return nil, err
	}

	var jobs []job
	for _, j := range res.Jobs {
		jobs = append(jobs, job{
			ID:     j.JobId,
			Status: j.Status,
			Result: j.Result,
		})
	}
	return jobs, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.jobs)-1 {
				m.cursor++
			}
		case "r":
			jobs, err := fetchJobs()
			m.jobs = jobs
			m.err = err
			m.cursor = 0
		case "/":
			m.filter = "queued" // hardcoded demo: filter queued jobs
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	output := style.Render("Job Status Dashboard\n\n")
	filteredJobs := m.jobs
	if m.filter != "" {
		var temp []job
		for _, j := range m.jobs {
			if j.Status == m.filter {
				temp = append(temp, j)
			}
		}
		filteredJobs = temp
	}
	for i, j := range filteredJobs {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		output += fmt.Sprintf("%s %s | %s | %s\n", cursor, j.ID, j.Status, j.Result)
	}
	output += "\n[↑↓] Navigate  [r] Refresh  [/] Filter queued  [q] Quit"
	return output
}

func main() {
	if err := tea.NewProgram(initialModel()).Start(); err != nil {
		log.Fatal(err)
	}
}
