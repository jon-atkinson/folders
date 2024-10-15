package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Intern_GetFoldersByOrgID(t *testing.T) {
	t.Run("test intern's implementation GetAllChildFolders - happy path", func(t *testing.T) {
		expected := []folder.Folder{
			{
				Name:  "driven-stripperella",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.driven-stripperella",
			},
			{
				Name:  "endless-master-mold",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.driven-stripperella.endless-master-mold",
			},
			{
				Name:  "valid-mega-man",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.driven-stripperella.valid-mega-man",
			},
			{
				Name:  "stirred-judomaster",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.driven-stripperella.stirred-judomaster",
			},
			{
				Name:  "complete-lockjaw",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.driven-stripperella.complete-lockjaw",
			},
			{
				Name:  "valued-captain",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.valued-captain",
			},
			{
				Name:  "frank-thunder",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.valued-captain.frank-thunder",
			},
			{
				Name:  "polished-bella",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.valued-captain.polished-bella",
			},
			{
				Name:  "proper-grim-reaper",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.valued-captain.proper-grim-reaper",
			},
			{
				Name:  "adapted-timeslip",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.adapted-timeslip",
			},
			{
				Name:  "learning-unicorn",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.adapted-timeslip.learning-unicorn",
			},
			{
				Name:  "pretty-firefly",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.adapted-timeslip.pretty-firefly",
			},
			{
				Name:  "innocent-eradicator",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.adapted-timeslip.innocent-eradicator",
			},
			{
				Name:  "faithful-warstar",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.adapted-timeslip.faithful-warstar",
			},
			{
				Name:  "thorough-miracleman",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.thorough-miracleman",
			},
			{
				Name:  "outgoing-cobweb",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.thorough-miracleman.outgoing-cobweb",
			},
			{
				Name:  "novel-squirrel",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.thorough-miracleman.novel-squirrel",
			},
			{
				Name:  "awake-cable",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.thorough-miracleman.awake-cable",
			},
			{
				Name:  "aware-smiling-tiger",
				OrgId: uuid.FromStringOrNil(folder.DefaultOrgID),
				Paths: "noble-vixen.nearby-secret.hip-stingray.thorough-miracleman.aware-smiling-tiger",
			},
		}

		folders := folder.GetSampleData()

		f := folder.NewDriver(folders)

		orgID := uuid.FromStringOrNil(folder.DefaultOrgID)
		cf := f.GetAllChildFolders(orgID, "hip-stingray")

		folder.SortFoldersByPath(expected)
		folder.SortFoldersByPath(cf)

		assert.EqualValues(t, expected, cf)
	})

	t.Run("test intern's implementation GetFoldersByOrgID - leaf node", func(t *testing.T) {
		folders := folder.GetSampleData()

		f := folder.NewDriver(folders)

		orgID := uuid.FromStringOrNil(folder.DefaultOrgID)
		cf := f.GetAllChildFolders(orgID, "emerging-nova")

		// some students return error, some don't
		if cf == nil {
			assert.Nil(t, cf)
		} else {
			assert.EqualValues(t, []folder.Folder{}, cf)
		}
	})

	t.Run("test intern's implementation GetFoldersByOrgID - mismatch orgID", func(t *testing.T) {
		folders := folder.GetSampleData()

		f := folder.NewDriver(folders)

		orgID := uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7")
		cf := f.GetAllChildFolders(orgID, "hip-stingray")
		assert.Nil(t, cf)

		// some students return error, some don't
	})

	t.Run("test intern's implementation GetFoldersByOrgID - mismatch folder", func(t *testing.T) {
		folders := folder.GetSampleData()

		f := folder.NewDriver(folders)

		orgID := uuid.FromStringOrNil(folder.DefaultOrgID)
		cf := f.GetAllChildFolders(orgID, "central-the-anarchis")
		assert.Nil(t, cf)

		// some students return error, some don't
	})
}

func Test_Itern_MoveFolder(t *testing.T) {
	t.Run("test intern's implementation MoveFolder - happy path", func(t *testing.T) {
		folders := folder.GetSampleData()

		f := folder.NewDriver(folders)

		_, err := f.MoveFolder("sacred-moonstar", "nearby-secret")

		res := f.GetAllChildFolders(uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"), "sacred-moonstar")

		expected := []folder.Folder{
			{
				Name:  "loved-retro-girl",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl",
			},
			{
				Name:  "safe-infragirl",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.safe-infragirl",
			},
			{
				Name:  "sweeping-hulkling",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.safe-infragirl.sweeping-hulkling",
			},
			{
				Name:  "elegant-silver-sable",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.safe-infragirl.elegant-silver-sable",
			},
			{
				Name:  "settling-blink",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.safe-infragirl.settling-blink",
			},
			{
				Name:  "worthy-cybergirl",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.worthy-cybergirl",
			},
			{
				Name:  "gentle-killmonger",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.worthy-cybergirl.gentle-killmonger",
			},
			{
				Name:  "helped-ultrawoman",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.worthy-cybergirl.helped-ultrawoman",
			},
			{
				Name:  "composed-wallflower",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower",
			},
			{
				Name:  "mutual-jigsaw",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.mutual-jigsaw",
			},
			{
				Name:  "measured-morbius",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.mutual-jigsaw.measured-morbius",
			},
			{
				Name:  "peaceful-metal-master",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.mutual-jigsaw.peaceful-metal-master",
			},
			{
				Name:  "moving-bizarro",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.moving-bizarro",
			},
			{
				Name:  "mature-coagula",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.moving-bizarro.mature-coagula",
			},
			{
				Name:  "positive-sentry",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.moving-bizarro.positive-sentry",
			},
			{
				Name:  "tight-titaness",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.tight-titaness",
			},
			{
				Name:  "novel-lettuce",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.tight-titaness.novel-lettuce",
			},
			{
				Name:  "sharp-glitter",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.tight-titaness.sharp-glitter",
			},
			{
				Name:  "unique-cherry",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.unique-cherry",
			},
			{
				Name:  "calm-penguin",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.unique-cherry.calm-penguin",
			},
			{
				Name:  "nearby-maestro",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro",
			},
			{
				Name:  "picked-glory",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.picked-glory",
			},
			{
				Name:  "gorgeous-wasp",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.picked-glory.gorgeous-wasp",
			},
			{
				Name:  "first-dragon-man",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.picked-glory.first-dragon-man",
			},
			{
				Name:  "mature-slipstream",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.picked-glory.mature-slipstream",
			},
			{
				Name:  "star-stormtrooper",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.picked-glory.star-stormtrooper",
			},
			{
				Name:  "dashing-forearm",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.dashing-forearm",
			},
			{
				Name:  "clear-supergran",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.dashing-forearm.clear-supergran",
			},
			{
				Name:  "related-kitty",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.dashing-forearm.related-kitty",
			},
			{
				Name:  "organic-hulk",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.dashing-forearm.organic-hulk",
			},
			{
				Name:  "healthy-deathstrike",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.healthy-deathstrike",
			},
			{
				Name:  "better-rapture",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.healthy-deathstrike.better-rapture",
			},
			{
				Name:  "enabled-professor-monster",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.enabled-professor-monster",
			},
			{
				Name:  "glowing-elongated",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.enabled-professor-monster.glowing-elongated",
			},
			{
				Name:  "equipped-hypno-hustler",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.enabled-professor-monster.equipped-hypno-hustler",
			},
		}

		folder.SortFoldersByPath(expected)
		folder.SortFoldersByPath(res)

		assert.NoError(t, err)
		assert.EqualValues(t, expected, res)
	})

	t.Run("test intern's implementation MoveFolder - multi move", func(t *testing.T) {
		folders := folder.GetSampleData()

		f := folder.NewDriver(folders)

		_, err := f.MoveFolder("sacred-moonstar", "nearby-secret")

		_, err = f.MoveFolder("dashing-forearm", "pure-blastaar")

		folder1 := f.GetAllChildFolders(uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"), "sacred-moonstar")

		folder2 := f.GetAllChildFolders(uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"), "dashing-forearm")

		expectedFolder1 := []folder.Folder{
			{
				Name:  "loved-retro-girl",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl",
			},
			{
				Name:  "safe-infragirl",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.safe-infragirl",
			},
			{
				Name:  "sweeping-hulkling",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.safe-infragirl.sweeping-hulkling",
			},
			{
				Name:  "elegant-silver-sable",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.safe-infragirl.elegant-silver-sable",
			},
			{
				Name:  "settling-blink",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.safe-infragirl.settling-blink",
			},
			{
				Name:  "worthy-cybergirl",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.worthy-cybergirl",
			},
			{
				Name:  "gentle-killmonger",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.worthy-cybergirl.gentle-killmonger",
			},
			{
				Name:  "helped-ultrawoman",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.loved-retro-girl.worthy-cybergirl.helped-ultrawoman",
			},
			{
				Name:  "composed-wallflower",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower",
			},
			{
				Name:  "mutual-jigsaw",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.mutual-jigsaw",
			},
			{
				Name:  "measured-morbius",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.mutual-jigsaw.measured-morbius",
			},
			{
				Name:  "peaceful-metal-master",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.mutual-jigsaw.peaceful-metal-master",
			},
			{
				Name:  "moving-bizarro",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.moving-bizarro",
			},
			{
				Name:  "mature-coagula",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.moving-bizarro.mature-coagula",
			},
			{
				Name:  "positive-sentry",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.moving-bizarro.positive-sentry",
			},
			{
				Name:  "tight-titaness",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.tight-titaness",
			},
			{
				Name:  "novel-lettuce",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.tight-titaness.novel-lettuce",
			},
			{
				Name:  "sharp-glitter",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.tight-titaness.sharp-glitter",
			},
			{
				Name:  "unique-cherry",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.unique-cherry",
			},
			{
				Name:  "calm-penguin",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.composed-wallflower.unique-cherry.calm-penguin",
			},
			{
				Name:  "nearby-maestro",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro",
			},
			{
				Name:  "picked-glory",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.picked-glory",
			},
			{
				Name:  "gorgeous-wasp",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.picked-glory.gorgeous-wasp",
			},
			{
				Name:  "first-dragon-man",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.picked-glory.first-dragon-man",
			},
			{
				Name:  "mature-slipstream",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.picked-glory.mature-slipstream",
			},
			{
				Name:  "star-stormtrooper",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.picked-glory.star-stormtrooper",
			},
			{
				Name:  "healthy-deathstrike",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.healthy-deathstrike",
			},
			{
				Name:  "better-rapture",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.healthy-deathstrike.better-rapture",
			},
			{
				Name:  "enabled-professor-monster",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.enabled-professor-monster",
			},
			{
				Name:  "glowing-elongated",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.enabled-professor-monster.glowing-elongated",
			},
			{
				Name:  "equipped-hypno-hustler",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "noble-vixen.nearby-secret.sacred-moonstar.nearby-maestro.enabled-professor-monster.equipped-hypno-hustler",
			},
		}

		expectedFolder2 := []folder.Folder{
			{
				Name:  "clear-supergran",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "stunning-horridus.pure-blastaar.dashing-forearm.clear-supergran",
			},
			{
				Name:  "related-kitty",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "stunning-horridus.pure-blastaar.dashing-forearm.related-kitty",
			},
			{
				Name:  "organic-hulk",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "stunning-horridus.pure-blastaar.dashing-forearm.organic-hulk",
			},
		}

		folder.SortFoldersByPath(expectedFolder1)
		folder.SortFoldersByPath(folder1)
		folder.SortFoldersByPath(expectedFolder2)
		folder.SortFoldersByPath(folder2)

		assert.NoError(t, err)
		assert.EqualValues(t, expectedFolder1, folder1)
		assert.EqualValues(t, expectedFolder2, folder2)
	})

	t.Run("test intern's implementation MoveFolder - leaf folder to leaf folder", func(t *testing.T) {
		folders := folder.GetSampleData()

		f := folder.NewDriver(folders)

		_, err := f.MoveFolder("related-kitty", "organic-hulk")

		cf := f.GetAllChildFolders(uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"), "dashing-forearm")

		expected := []folder.Folder{
			{
				Name:  "clear-supergran",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "stunning-horridus.sacred-moonstar.nearby-maestro.dashing-forearm.clear-supergran",
			},
			{
				Name:  "related-kitty",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "stunning-horridus.sacred-moonstar.nearby-maestro.dashing-forearm.organic-hulk.related-kitty",
			},
			{
				Name:  "organic-hulk",
				OrgId: uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17a"),
				Paths: "stunning-horridus.sacred-moonstar.nearby-maestro.dashing-forearm.organic-hulk",
			},
		}

		folder.SortFoldersByPath(expected)
		folder.SortFoldersByPath(cf)

		assert.NoError(t, err)
		assert.EqualValues(t, expected, cf)
	})

	t.Run("test intern's implementation MoveFolder - invalid source path", func(t *testing.T) {
		folders := folder.GetSampleData()
		f := folder.NewDriver(folders)

		_, err := f.MoveFolder("weird-source", "nearby-maestro")

		assert.Error(t, err)
	})

	t.Run("test intern's implementation MoveFolder - invalid destination path", func(t *testing.T) {
		folders := folder.GetSampleData()
		f := folder.NewDriver(folders)

		_, err := f.MoveFolder("nearby-maestro", "weird-destination")

		assert.Error(t, err)
	})

	t.Run("test intern's implementation MoveFolder - cross org folder movement", func(t *testing.T) {
		folders := folder.GetSampleData()
		f := folder.NewDriver(folders)

		_, err := f.MoveFolder("sacred-moonstar", "steady-insect")

		assert.Error(t, err)
	})

	t.Run("test intern's implementation MoveFolder - move into itself", func(t *testing.T) {
		folders := folder.GetSampleData()
		f := folder.NewDriver(folders)

		_, err := f.MoveFolder("sacred-moonstar", "sacred-moonstar")

		assert.Error(t, err)
	})

	t.Run("test intern's implementation MoveFolder - move into child folder", func(t *testing.T) {
		folders := folder.GetSampleData()
		f := folder.NewDriver(folders)

		_, err := f.MoveFolder("sacred-moonstar", "elegant-silver-sable")

		assert.Error(t, err)
	})
}
