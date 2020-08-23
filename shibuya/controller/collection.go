package controller

import (
	"github.com/harpratap/shibuya/shibuya/model"
	"github.com/harpratap/shibuya/shibuya/utils"
)

func prepareCollection(collection *model.Collection) []*EngineDataConfig {
	planCount := len(collection.ExecutionPlans)
	edc := EngineDataConfig{
		EngineData: map[string]*model.ShibuyaFile{},
	}
	engineDataConfigs := edc.deepCopies(planCount)
	for i := 0; i < planCount; i++ {
		for _, d := range collection.Data {
			sf := model.ShibuyaFile{
				Filename:     d.Filename,
				Filepath:     d.Filepath,
				TotalSplits:  1,
				CurrentSplit: 0,
			}
			if collection.CSVSplit {
				sf.TotalSplits = planCount
				sf.CurrentSplit = i
			}
			engineDataConfigs[i].EngineData[sf.Filename] = &sf
		}
	}
	return engineDataConfigs
}

func (c *Controller) TermAndPurgeCollection(collection *model.Collection) error {
	// This is a force remove so we ignore the errors happened at test termination
	c.TermCollection(collection, true)
	err := utils.Retry(func() error {
		return c.Kcm.PurgeCollection(collection.ID)
	})
	return err
}
