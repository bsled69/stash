package autotag

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPerformerScenes(t *testing.T) {
	type test struct {
		performerName string
		expectedRegex string
	}

	performerNames := []test{
		{
			"performer name",
			`(?i)(?:^|_|[^\w\d])performer[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
		{
			"performer + name",
			`(?i)(?:^|_|[^\w\d])performer[.\-_ ]*\+[.\-_ ]*name(?:$|_|[^\w\d])`,
		},
	}

	for _, p := range performerNames {
		testPerformerScenes(t, p.performerName, p.expectedRegex)
	}
}

func testPerformerScenes(t *testing.T, performerName, expectedRegex string) {
	mockSceneReader := &mocks.SceneReaderWriter{}

	const performerID = 2

	var scenes []*models.Scene
	matchingPaths, falsePaths := generateScenePaths(performerName)
	for i, p := range append(matchingPaths, falsePaths...) {
		scenes = append(scenes, &models.Scene{
			ID:   i + 1,
			Path: p,
		})
	}

	performer := models.Performer{
		ID:   performerID,
		Name: models.NullString(performerName),
	}

	organized := false
	perPage := models.PerPageAll

	expectedSceneFilter := &models.SceneFilterType{
		Organized: &organized,
		Path: &models.StringCriterionInput{
			Value:    expectedRegex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
	}

	expectedFindFilter := &models.FindFilterType{
		PerPage: &perPage,
	}

	mockSceneReader.On("Query", expectedSceneFilter, expectedFindFilter).Return(scenes, len(scenes), nil).Once()

	for i := range matchingPaths {
		sceneID := i + 1
		mockSceneReader.On("GetPerformerIDs", sceneID).Return(nil, nil).Once()
		mockSceneReader.On("UpdatePerformers", sceneID, []int{performerID}).Return(nil).Once()
	}

	err := PerformerScenes(&performer, nil, mockSceneReader)

	assert := assert.New(t)

	assert.Nil(err)
	mockSceneReader.AssertExpectations(t)
}