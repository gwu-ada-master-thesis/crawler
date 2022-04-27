package master

import (
	// built-in
	"log"
	"sync"
	"time"

	// 3rd party

	// custom
	"masterthesis/utils"
)

func (master *Master) RequestJob(args *JobRequest, reply *JobResponse) error {

	// print master status
	log.Printf("processing task request\n")

	// lock file to single worker
	master.Access.Lock()

	// unlock file for other workers
	defer master.Access.Unlock()

	for _, job := range master.Jobs {

		// in case an idle job is found
		if job.Status == "idle" && master.ProcessDetails.domainMetadataApiRequestCount < 45 {

			// update job status
			job.Status = "occupied"

			// send response
			reply.Job = job
			reply.Status = "assigned"

			// update job details
			master.ProcessDetails.idleJobsCount = master.ProcessDetails.idleJobsCount - 1
			master.ProcessDetails.occupiedJobsCount = master.ProcessDetails.occupiedJobsCount + 1

			// print master status
			log.Printf("assigning -> task = %s, status = %s\n", job.Task, job.Status)

			// // monitor worker activity
			go master.monitorWorker(job)

			return nil

		}

	}

	return nil
}

func (master *Master) ReportJob(args *JobRequest, reply *JobResponse) error {

	// print master status
	log.Printf("documenting worker report -> task = %s\n", args.Job.Task)

	// lock file to single worker
	master.Access.Lock()

	// unlock file for other workers
	defer master.Access.Unlock()

	// get job
	job := master.Jobs[args.Job.Task]

	// print master status
	log.Printf("updating -> task = %s, status = %s to	done\n", job.Task, job.Status)

	// update job status
	job.Status = "done"

	// update job details
	master.ProcessDetails.occupiedJobsCount = master.ProcessDetails.occupiedJobsCount - 1
	master.ProcessDetails.doneJobsCount = master.ProcessDetails.doneJobsCount + 1

	// update callback status
	reply.Status = "success"

	return nil
}

func (master *Master) monitorWorker(job *Job) {

	// initiate timer for 10 second
	timer := time.NewTicker(1 * time.Hour)
	defer timer.Stop()

	// infinitely loop
	for {

		// and if
		select {

		// in case, timer is up
		case <-timer.C:

			// lock file to single worker
			master.Access.Lock()

			// revoke job
			job.Status = "idle"

			// update job details
			master.ProcessDetails.occupiedJobsCount = master.ProcessDetails.occupiedJobsCount - 1
			master.ProcessDetails.idleJobsCount = master.ProcessDetails.idleJobsCount + 1

			// unlock file for other workers
			master.Access.Unlock()

			// end monitoring
			return

		// otherwise
		default:

			// if job is done
			if job.Status == "done" {

				// end monitoring
				return

			}

		}
	}
}

func (master *Master) monitorResources() {

	// initiate timer for a minute
	timer := time.NewTicker(60 * time.Second)

	defer timer.Stop()

	// infinitely loop
	for {

		// and if
		select {

		// in case, timer is up
		case <-timer.C:

			log.Printf("resetting metadata api request count\n")

			// reset resource
			master.ProcessDetails.domainMetadataApiRequestCount = 0

			// end monitoring
			return

		// otherwise
		default:
			continue
		}
	}
}

func (master *Master) Done() bool {

	// reserve access
	master.Access.Lock()

	// release access
	defer master.Access.Unlock()

	log.Printf("showing job statuses: idle job count - %d, occupied job count - %d, done - %d\n", master.ProcessDetails.idleJobsCount, master.ProcessDetails.occupiedJobsCount, master.ProcessDetails.doneJobsCount)

	// in case all tasks are done
	return (master.ProcessDetails.idleJobsCount + master.ProcessDetails.occupiedJobsCount) == 0
}

func (master *Master) addJobs(tasks []string) {

	// reserve access
	master.Access.Lock()

	// release access
	defer master.Access.Unlock()

	// loop each task
	for _, task := range tasks {

		id := master.ProcessDetails.latestJobId

		// add task to master's jobs
		master.Jobs[task] = &Job{Id: id, Task: task, Status: "idle"}

		// update latest job id
		master.ProcessDetails.latestJobId = master.ProcessDetails.latestJobId + 1
	}

	// update idle jobs count
	master.ProcessDetails.idleJobsCount += len(tasks)
}

func initialize() *ProcessDetails {

	// load environment
	utils.LoadEnvironment()

	// create initial process values
	processDetails := ProcessDetails{
		latestJobId: 1,

		idleJobsCount:     0,
		occupiedJobsCount: 0,
		doneJobsCount:     0,

		domainMetadataApiRequestCount: 0,
	}

	return &processDetails
}

func MakeMaster(tasks []string) *Master {

	//
	// ----------------------------------------------------------------* variables
	//

	// initialize
	processDetails := initialize()

	// initialize master
	master := Master{
		Jobs:           make(map[string]*Job),
		Access:         new(sync.RWMutex),
		SocketName:     utils.GenerateSocketName(824, "master"),
		ProcessDetails: processDetails,
	}

	//
	// ----------------------------------------------------------------* processes
	//

	// print master status
	log.Printf("adding jobs\n")
	time.Sleep(1 * time.Second)

	// add inital jobs
	master.addJobs(tasks)

	// print master status
	log.Printf("starting resource monitoring\n")
	time.Sleep(1 * time.Second)

	// monitor resources in a spare thread
	go master.monitorResources()

	// print master status
	log.Printf("starting server\n")
	time.Sleep(1 * time.Second)

	// start server
	master.serve()

	return &master
}
