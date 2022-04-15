package main

import (
	"embed"
	"fmt"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xuri/excelize/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"image/png"
	"log"
	"strconv"
)

//go:embed graphics/*
var EmbeddedAssets embed.FS

var counter = 0
var demoApp GuiApp
var textWidget *widget.Text

func main() {
	ebiten.SetWindowSize(700, 600)
	ebiten.SetWindowTitle("Project 4")

	demoApp = GuiApp{AppUI: MakeUIWindow()}

	err := ebiten.RunGame(&demoApp)
	if err != nil {
		log.Fatalln("Error running User Interface Demo", err)
	}
}

func (g GuiApp) Update() error {
	//TODO finish me
	g.AppUI.Update()
	return nil
}

func (g GuiApp) Draw(screen *ebiten.Image) {
	//TODO finish me
	g.AppUI.Draw(screen)
}

func (g GuiApp) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

type GuiApp struct {
	AppUI *ebitenui.UI
}

func MakeUIWindow() (GUIhandler *ebitenui.UI) {
	background := image.NewNineSliceColor(color.Gray16{})
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(background))
	textInfo := widget.TextOptions{}.Text("This is our first Window", basicfont.Face7x13, color.White)

	idle, err := loadImageNineSlice("button-idle.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	hover, err := loadImageNineSlice("button-hover.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	pressed, err := loadImageNineSlice("button-pressed.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	disabled, err := loadImageNineSlice("button-disabled.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	buttonImage := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
	button := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Press Me", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		// ... click handler, etc. ...
		widget.ButtonOpts.ClickedHandler(FunctionNameHere),
	)
	rootContainer.AddChild(button)
	resources, err := newListResources()
	if err != nil {
		log.Println(err)
	}

	allStates := loadStates()
	dataAsGeneric := make([]interface{}, len(allStates))
	for position, states := range allStates {
		dataAsGeneric[position] = states
	}

	listWidget := widget.NewList(
		widget.ListOpts.Entries(dataAsGeneric),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			stateInformation := fmt.Sprintf("%s %s %s", e.(States).StateName, e.(States).NPOPCHG2020, e.(States).NPOPCHG2021)
			return stateInformation
		}),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(resources.image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(resources.track, resources.handle),
			widget.SliderOpts.HandleSize(resources.handleSize),
			widget.SliderOpts.TrackPadding(resources.trackPadding)),
		widget.ListOpts.EntryColor(resources.entry),
		widget.ListOpts.EntryFontFace(resources.face),
		widget.ListOpts.EntryTextPadding(resources.entryPadding),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			percent := args.Entry.(States).PERCENTCHANGE
			textWidget.Label = percent
			//widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			//stateInformation := fmt.Sprintf("%s", e.(States).PERCENTCHANGE)
			//textWidget.Label = stateInformation
			//return stateInformation
			//want to use args.Entry to identify location of list entry and use it to index into struct but how do I do this
			//above code will not work because even if it did would just print that entire field from the struct at once instead
			//of wanted index
			//do I need to make another slice and index based of that with args.Entry args.Entry has to equal something
			//args.Entry is an interface but, I would like an int from the location for it
			//args.entry is selected and cast entry to member of struct
			//EntryLabelFunc needs a return statement
			//do something when a list item changes
			//args. and see what happens and may have to use e interface
			//})
		}))
	rootContainer.AddChild(listWidget)
	textWidget = widget.NewText(textInfo)
	rootContainer.AddChild(textWidget)

	GUIhandler = &ebitenui.UI{Container: rootContainer}
	return GUIhandler
}

func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	i := loadPNGImageFromEmbedded(path)

	w, h := i.Size()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("graphics")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("graphics/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func FunctionNameHere(args *widget.ButtonClickedEventArgs) {
	counter++
	message := fmt.Sprintf("I've gotten coffee this semester %d times", counter)
	textWidget.Label = message
}

//https://www.socketloop.com/tutorials/golang-calculate-percentage-change-of-two-values
func PercentageChange(old, new int) (delta float64) {
	diff := float64(new - old)
	delta = (diff / float64(old)) * 100
	return
}

func loadStates() []States {
	listOfStates := make([]States, 51)
	stateLocation := 0
	excelFile, err := excelize.OpenFile("countyPopChange2020-2021.xlsx")
	if err != nil {
		log.Fatalln(err)
	}
	all_rows, err := excelFile.GetRows("co-est2021-alldata")
	if err != nil {
		log.Fatalln(err)
	}
	for number, row := range all_rows {
		if number < 1 {
			continue
		}
		if len(row) <= 1 {
			continue
		}
		if row[4] == "0" {
			new2021, _ := strconv.Atoi(row[11])
			pop2021, _ := strconv.Atoi(row[9])
			currentState := States{
				StateName:     fmt.Sprintf("State : %s", row[5]),
				NPOPCHG2020:   fmt.Sprintf("Population Change 2020 : %s", row[10]),
				NPOPCHG2021:   fmt.Sprintf("Population Change 2021 : %s", row[11]),
				PERCENTCHANGE: fmt.Sprintf("Percent Change : %0.2f%% \n", PercentageChange(int(pop2021), int(new2021))),
			}
			listOfStates[stateLocation] = currentState
			stateLocation++
		}
	}
	return listOfStates
}
