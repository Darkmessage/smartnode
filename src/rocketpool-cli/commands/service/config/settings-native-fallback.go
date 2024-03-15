package config

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rocket-pool/node-manager-core/config/ids"
	"github.com/rocket-pool/smartnode/shared/config"
)

// The page wrapper for the fallback config
type NativeFallbackConfigPage struct {
	home           *settingsNativeHome
	page           *page
	layout         *standardLayout
	masterConfig   *config.SmartNodeConfig
	useFallbackBox *parameterizedFormItem
	reconnectDelay *parameterizedFormItem
	fallbackItems  []*parameterizedFormItem
}

// Creates a new page for the fallback client settings
func NewNativeFallbackConfigPage(home *settingsNativeHome) *NativeFallbackConfigPage {
	configPage := &NativeFallbackConfigPage{
		home:         home,
		masterConfig: home.md.Config,
	}
	configPage.createContent()

	configPage.page = newPage(
		home.homePage,
		"settings-native-fallback",
		"Fallback Clients",
		"Select this to specify a separate pair of externally-managed Execution Client and Beacon Node that the Smart Node will use if your main clients ever go offline.",
		configPage.layout.grid,
	)

	return configPage
}

// Get the underlying page
func (configPage *NativeFallbackConfigPage) getPage() *page {
	return configPage.page
}

// Creates the content for the fallback client settings page
func (configPage *NativeFallbackConfigPage) createContent() {

	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(&configPage.masterConfig.Network, "Fallback Client Settings")

	// Return to the home page after pressing Escape
	configPage.layout.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			// Close all dropdowns and break if one was open
			for _, param := range configPage.layout.parameters {
				dropDown, ok := param.item.(*DropDown)
				if ok && dropDown.open {
					dropDown.CloseList(configPage.home.md.app)
					return nil
				}
			}

			// Return to the home page
			configPage.home.md.setPage(configPage.home.homePage)
			return nil
		}
		return event
	})

	// Set up the form items
	configPage.useFallbackBox = createParameterizedCheckbox(&configPage.masterConfig.Fallback.UseFallbackClients)
	configPage.fallbackItems = createParameterizedFormItems(configPage.masterConfig.Fallback.GetParameters(), configPage.layout.descriptionBox)

	// Take the enable out since it's done explicitly
	fallbackItems := []*parameterizedFormItem{}
	for _, item := range configPage.fallbackItems {
		if item.parameter.GetCommon().ID == ids.FallbackUseFallbackClientsID {
			continue
		}
		fallbackItems = append(fallbackItems, item)
	}
	configPage.fallbackItems = fallbackItems

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.useFallbackBox)
	configPage.layout.mapParameterizedFormItems(configPage.fallbackItems...)

	// Set up the setting callbacks
	configPage.useFallbackBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if configPage.masterConfig.Fallback.UseFallbackClients.Value == checked {
			return
		}
		configPage.masterConfig.Fallback.UseFallbackClients.Value = checked
		configPage.handleUseFallbackChanged()
	})

	// Do the initial draw
	configPage.handleUseFallbackChanged()
}

// Handle all of the form changes when the Use Fallback box has changed
func (configPage *NativeFallbackConfigPage) handleUseFallbackChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.useFallbackBox.item)

	// Only add the supporting stuff if external clients are enabled
	if configPage.masterConfig.Fallback.UseFallbackClients.Value == false {
		return
	}
	configPage.layout.form.AddFormItem(configPage.reconnectDelay.item)
	configPage.layout.addFormItems(configPage.fallbackItems)

	configPage.layout.refresh()
}

// Handle a bulk redraw request
func (configPage *NativeFallbackConfigPage) handleLayoutChanged() {
	configPage.handleUseFallbackChanged()
}
