// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Home() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\" data-theme=\"synthwave\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Your Coffee is Brewferring</title><link rel=\"stylesheet\" href=\"/static/css/output.css\"><script src=\"/static/script/htmx.min.js\"></script></head><body class=\"min-h-screen bg-base-200\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Navbar(NavbarProps{IsAuthenticated: false}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<main class=\"container mx-auto px-4 py-8\"><div class=\"hero min-h-[60vh] bg-base-200 rounded-box\"><div class=\"hero-content text-center\"><div class=\"max-w-md\"><h1 class=\"text-5xl font-bold text-primary\">Your Coffee is Brewferring</h1><p class=\"py-6\">Welcome to the finest coffee ordering experience. We've partnered with Terminal.shop to bring you the best coffee beans from around the world.</p><a href=\"/products\" class=\"btn btn-primary\">View Products</a></div></div></div><div class=\"grid grid-cols-1 md:grid-cols-3 gap-6 mt-8\"><div class=\"card bg-base-100 shadow-xl\"><div class=\"card-body\"><h2 class=\"card-title\">Premium Beans</h2><p>Carefully selected coffee beans from the world's best regions.</p></div></div><div class=\"card bg-base-100 shadow-xl\"><div class=\"card-body\"><h2 class=\"card-title\">Fresh Roasting</h2><p>Each batch is freshly roasted to bring out the perfect flavor profile.</p></div></div><div class=\"card bg-base-100 shadow-xl\"><div class=\"card-body\"><h2 class=\"card-title\">Fast Delivery</h2><p>Quick and reliable shipping to get your coffee to you fresh.</p></div></div></div></main></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
