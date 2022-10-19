package schema

const (
	QueryParam  ParamLocation = "query"
	HeaderParam ParamLocation = "header"
	PathParam   ParamLocation = "path"
	CookieParam ParamLocation = "cookie"
)

type API struct {
	Info    Info            `json:"info"`
	Servers []Server        `json:"servers"`
	Paths   map[string]Path `json:"paths"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Path struct {
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	Get         *Operation `json:"get,omitempty"`
	Put         *Operation `json:"put,omitempty"`
	Post        *Operation `json:"post,omitempty"`
	Delete      *Operation `json:"delete,omitempty"`
	Options     *Operation `json:"options,omitempty"`
	Head        *Operation `json:"head,omitempty"`
	Patch       *Operation `json:"patch,omitempty"`
	Trace       *Operation `json:"trace,omitempty"`
}

type Operation struct {
	Params   []Param             `json:"parameters,omitempty"`
	ReqBody  ReqBody             `json:"requestBody"`
	Response map[string]Response `json:"response"`
}

type Param struct {
	Name        string        `json:"name"`
	In          ParamLocation `json:"in"`
	Description string        `json:"description,omitempty"`
	Required    bool          `json:"required,omitempty"`
}

type ParamLocation string

type ReqBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
	Required    bool                 `json:"required,omitempty"`
}

type Response struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
}

type MediaType struct {
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}
