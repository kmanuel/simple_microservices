package api_task

import "strconv"

func mapCropTaskToFaktoryTask(task *CropTask) *Task {
	return &Task{
		ID: task.ID,
		Type: "crop",
		TaskParams: map[string]interface{}{
			"image_id": task.ImageId,
			"height": strconv.Itoa(task.Height),
			"width": strconv.Itoa(task.Width),
		},
	}

}

func mapScreenshotTaskToFaktoryTask(screenshotTask *ScreenShotTask) *Task {
	return &Task{
		ID: screenshotTask.ID,
		Type: "screenshot",
		TaskParams: map[string]interface{}{
			"url": screenshotTask.Url,
		},
	}
}
