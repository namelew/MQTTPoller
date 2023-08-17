package experiments

import (
	"time"

	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/databases"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/databases/models"
	"github.com/namelew/mqtt-bm-latency/packages/logs"
)

type Experiments struct {
	log *logs.Log
}

const id = 2

func Build(l *logs.Log) *Experiments {
	return &Experiments{
		log: l,
	}
}

func (h *Experiments) Add(e models.Experiment, d models.ExperimentDeclaration, wid ...int) error {
	var wkrs []models.Worker
	var declaration, empty models.ExperimentDeclaration

	if len(wid) <= 1 || wid[0] == -1 {
		if err := (databases.DB.Find(&wkrs)).Error; err != nil {
			h.log.Register("Unable to find workers")
			return err
		}
	} else {
		if err := (databases.DB.Where("id IN ?", wid).Find(&wkrs)).Error; err != nil {
			h.log.Register("Unable to find workers")
			return err
		}
	}

	if err := (databases.DB.Model(&d).Find(&declaration)).Error; err != nil {
		h.log.Register("Unable to query experiment declaration")
		return err
	}

	if declaration == empty {
		if err := (databases.DB.Create(&d)).Error; err != nil {
			h.log.Register("Unable to create experiment declaration")
			return err
		}
	}

	for i := range wkrs {
		e.Workers = append(e.Workers, &wkrs[i])
	}

	e.ExperimentDeclarationID = declaration.ID
	e.ExperimentDeclaration = declaration

	if err := (databases.DB.Create(&e)).Error; err != nil {
		h.log.Register("Unable to register experiment")
		return err
	}

	return nil
}

func (h *Experiments) Remove(id uint64) error {
	if err := (databases.DB.Model(&models.Experiment{}).Where("id = ?", id).Delete(&models.Experiment{})).Error; err != nil {
		h.log.Fatal("Unable to remove experiment")
		return err
	}

	return nil
}

func (h *Experiments) Update(key uint64, new models.Experiment) error {
	var experiment models.Experiment

	if err := (databases.DB.Model(&models.Experiment{}).Where("id = ?", key).Find(&experiment)).Error; err != nil || experiment.ID == 0 {
		h.log.Register("Unable to find selected experiment")
		return err
	}

	if err := (databases.DB.Save(&experiment)).Error; err != nil {
		h.log.Register("Unable to update experiment")
		return err
	}

	return nil
}

func (h *Experiments) List() ([]models.Experiment, error) {
	var experiments []models.Experiment

	if err := (databases.DB.Find(&experiments)).Error; err != nil {
		h.log.Register("Unable to get experiments registers on database")
		return experiments, err
	}

	return experiments, nil
}

func (h *Experiments) Get(id uint64) (models.Experiment, error) {
	var experiment models.Experiment

	if err := (databases.DB.Where("id = ?", id).Find(&experiment)).Error; err != nil {
		h.log.Register("Unable to get experiment register on database")
		return experiment, err
	}

	return experiment, nil
}

func (h *Experiments) TrashOut(i time.Time) {
	if (databases.DB.Model(&models.Experiment{}).Where("finish = ? AND (created_at < ? OR updated_at < ?)", true, i, i).Delete(&models.Experiment{})).Error != nil {
		h.log.Fatal("Unable to remove finished experiment")
	}
}

func (h *Experiments) ID() int {
	return id
}
