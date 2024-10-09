package entrypoint

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/script"
)

// Parse D:HH:MM:SS, HH:MM:SS, or MM:SS format into time.Duration
func parseRemainingTime(remainingTimeStr string) (time.Duration, error) {
	// Split the time string by ":"
	parts := strings.Split(remainingTimeStr, ":")
	var days, hours, minutes, seconds int

	if len(parts) == 4 {
		days, _ = strconv.Atoi(parts[0])
		hours, _ = strconv.Atoi(parts[1])
		minutes, _ = strconv.Atoi(parts[2])
		seconds, _ = strconv.Atoi(parts[3])
	} else if len(parts) == 3 {
		hours, _ = strconv.Atoi(parts[0])
		minutes, _ = strconv.Atoi(parts[1])
		seconds, _ = strconv.Atoi(parts[2])
	} else if len(parts) == 2 {
		minutes, _ = strconv.Atoi(parts[0])
		seconds, _ = strconv.Atoi(parts[1])
	} else {
		return 0, fmt.Errorf("invalid remaining time format: %s", remainingTimeStr)
	}

	// Calculate total duration in seconds
	totalDuration := time.Duration(days)*24*time.Hour + time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second

	return totalDuration, nil
}

func getSlurmEndTime() (time.Time, error) {
	// Get the SLURM job ID from the environment
	jobID := os.Getenv("SLURM_JOB_ID")
	if jobID == "" {
		fmt.Println("Not running within a SLURM job.")
		return time.Time{}, nil
	}

	// Prepare the squeue command to get the remaining time (%L)
	cmd := exec.Command("squeue", "-h", "-j", jobID, "-o", "%L")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing squeue:", err)
		return time.Time{}, err
	}

	remainingTimeStr := strings.TrimSpace(string(output))

	// Parse the remaining time (HH:MM:SS) into a time.Duration
	remainingTime, err := parseRemainingTime(remainingTimeStr)
	if err != nil {
		return time.Time{}, err
	}

	// Calculate the end time by adding the remaining time to the current time
	endTime := time.Now().Add(remainingTime)
	return endTime, nil
}

// Example usage
func setEndTimerIfSlurm() {
	// Check if running in SLURM
	if os.Getenv("SLURM_JOB_ID") != "" {
		getSlurmMetadata()
		endTime, err := getSlurmEndTime()
		if err != nil {
			log.Log.Warn("Error getting SLURM end time: %v", err)
			return
		}

		// Start a goroutine to notify when there are 15 seconds left
		for {
			remaining := time.Until(endTime)
			if remaining <= 30*time.Second && remaining > 0 {
				// If 30 seconds or less are remaining, print the message
				log.Log.Warn("30 seconds remaining until the job ends!")
				log.Log.Warn("Cleanly exiting the simulation early...")
				engine.Exit()
			}
			// Sleep for a short while before checking again
			time.Sleep(15 * time.Second)
		}
	}
}

func getSlurmMetadata() {
	script.MMetadata.Add("slurm_user", os.Getenv("SLURM_JOB_USER"))
	script.MMetadata.Add("slurm_partition", os.Getenv("SLURM_JOB_PARTITION"))
	script.MMetadata.Add("slurm_job_id", os.Getenv("SLURM_JOB_ID"))
	script.MMetadata.Add("slurm_node", os.Getenv("SLURM_JOB_NODELIST"))
}
