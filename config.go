package goui

import (
	"sync"
	"path/filepath"
	"flag"
	"strings"
	"log"
	"html/template"
	"github.com/spf13/viper"
	"runtime"
)

const (
	CfgVerbose = "verbose"
	CfgTemplatePath = "templatePaths"
	CfgHomepage = "homepage"
	CfgReload = "dynamicreload"
	CfgPattern = "tmplpattern"
)

type UIContext struct {
	sync.Mutex
	t *template.Template
	p *viper.Viper
}

var (
	defaultCfg UIContext
)

func currentPath() string {
	_, b, _, _ := runtime.Caller(0)
	dir := filepath.Dir(b)
	return dir
}

func init() {
	/* https://github.com/spf13/viper */
	defaultCfg.p = viper.New()
	defaultCfg.p.SetConfigType("json")
	defaultCfg.p.SetConfigName("goui.config.json") // name of config file (without extension)
	defaultCfg.p.AddConfigPath(currentPath())              // path to look for the config file in

	// Process command line flags, and set defaults
	verbose := *flag.Bool("verbose", false, "Turn on/off verbose log messages")
	reload := *flag.Bool("dynamicreload", false, "Reload (parse) templates for every request.")
	homepage := *flag.String("home", "index.html", "The default home page, e.g. index.html.")
	pattern := *flag.String("pattern", "*.html",
		"Filename pattern for loading templates. E.g. -pattern *.html.")
	uPathList := flag.String("templatePaths", strings.Join([]string{".", "templates"}, ";"),
		"Search paths for templates. Multiple paths are separated by semi-colons.")
	flag.Parse()

	// Find and read the config file
	err := defaultCfg.p.ReadInConfig()
	if err != nil {
		log.Printf("init(): %s \n", err)
	}

	// Watching and re-reading config files
	// Uses fsnotify/fsnotify for file event notications
	viper.WatchConfig()

	defaultCfg.p.Set(CfgVerbose, verbose)
	defaultCfg.p.Set(CfgReload, reload)
	defaultCfg.p.Set(CfgHomepage, homepage)
	defaultCfg.p.Set(CfgPattern, pattern)
	defaultCfg.p.Set(CfgTemplatePath, strings.Split(*uPathList, ";"))

	// Initialize the root/home page
	defaultCfg.t = template.New(viper.GetString(CfgHomepage))
	// Initialize the templates
	if paths := defaultCfg.p.GetStringSlice(CfgTemplatePath); len(paths) > 0 {
		appendTemplatePaths(&defaultCfg, paths)
	}
	defaultCfg.t.Funcs(template.FuncMap{
		"StripWhitespace": StripWhitespace,
	})

	log.Println("Default Context Settings:")
	for k, v := range defaultCfg.p.AllKeys() {
		log.Printf("%d) %s: %v\n", k, v, defaultCfg.p.Get(v))
	}

}

// Get the configuration information
func GetUIConfig() *UIContext {
	c := UIContext(defaultCfg)
	return &c
}


func NewUIContext() *UIContext {
	c := UIContext(defaultCfg)
	return &c
}

// Return the list of paths used to find templates.
func (uic *UIContext) TemplatePaths() []string {
	return uic.p.GetStringSlice(CfgTemplatePath)
}

// Add search paths to the configuration
func (uic *UIContext) AddTemplatePaths(paths...string) *UIContext {
	uic.Lock()
	defer uic.Unlock()
	appendTemplatePaths(uic, paths)
	return uic
}

// Load all templates from the search path, using the fPattern as the filename regex.
// If fPattern is the string zero-value, it used the context default.
func (uic *UIContext) LoadTemplates(fPattern string) {
	uic.Lock()
	defer uic.Unlock()
	if len(fPattern) == 0 {
		fPattern = uic.p.GetString(CfgPattern)
	}

	for _, v := range uic.p.GetStringSlice(CfgTemplatePath) {
		f := filepath.Join(v, fPattern)
		t, pErr := template.ParseGlob(f)
		if pErr != nil {
			log.Printf("LoadTemplates, %s.", pErr)
			continue
		}
		for _, tmpl := range t.Templates() {
			uic.t.AddParseTree(tmpl.Name(), tmpl.Tree)
		}
	}
}

// Get the template data for parse/execute
func (uic *UIContext) Templates() *template.Template {
	return uic.t
}

// Add paths to the template search list
// Assumes the caller performs locking/unlocking.
func appendTemplatePaths(uic *UIContext, paths []string) {
	var pathList = uic.p.GetStringSlice(CfgTemplatePath)
	for _, v := range paths {
		p, err := filepath.Abs(v)
		if err != nil {
			logMsg("appendTemplatePaths, invalid or non-existent path.", map[string]string{
				"error": err.Error(), "path": v})
			continue
		}
		pathList = append(pathList, p)
	}
	uic.p.Set(CfgTemplatePath, pathList)
}

func logMsg(msg string, params map[string]string) {

	log.Printf("%s: ", msg)
	for k, v := range params {
		log.Printf("\"%s\": %s ", k, v)
	}
	log.Println()
}
