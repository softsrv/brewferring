// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"context"
	"io"
	"strconv"

	"github.com/softsrv/brewferring/internal/models"
)

type SchedulersProps struct {
	Schedulers []models.Scheduler
	Devices    []models.Device
}

func Schedulers(schedulers []models.Scheduler, devices []models.Device) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		props := SchedulersProps{
			Schedulers: schedulers,
			Devices:    devices,
		}
		return SchedulersTemplate(props).Render(ctx, w)
	})
}

func SchedulersTemplate(props SchedulersProps) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Schedulers - Brewferring</title><link href=\"https://cdn.jsdelivr.net/npm/daisyui@2.6.0/dist/full.css\" rel=\"stylesheet\" type=\"text/css\"><script src=\"https://cdn.tailwindcss.com\"></script></head><body class=\"bg-base-200\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Navbar(NavbarProps{IsAuthenticated: true}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<div class=\"container mx-auto px-4 py-8\"><div class=\"flex justify-between items-center mb-8\"><h1 class=\"text-3xl font-bold\">Schedulers</h1><button class=\"btn btn-primary\" onclick=\"document.getElementById(&#39;create-scheduler-modal&#39;).showModal()\">Create Scheduler</button></div><div class=\"overflow-x-auto\"><table class=\"table w-full\"><thead><tr><th>Name</th><th>Type</th><th>Details</th><th>Actions</th></tr></thead> <tbody>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, scheduler := range props.Schedulers {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<tr><td>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(scheduler.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/schedulers.templ`, Line: 57, Col: 28}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</td><td>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(scheduler.Type)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/schedulers.templ`, Line: 58, Col: 28}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</td><td>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if scheduler.Type == "device" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "Device: ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(scheduler.Device.Name)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/schedulers.templ`, Line: 61, Col: 41}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, " (Threshold: ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var5 string
				templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.FormatFloat(scheduler.Threshold, 'f', 1, 64))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/schedulers.templ`, Line: 61, Col: 108}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, ")")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "Date: ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(scheduler.Date.Format("2006-01-02 15:04:05"))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/schedulers.templ`, Line: 63, Col: 62}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "</td><td><button class=\"btn btn-sm btn-error\" onclick=\"deleteScheduler({strconv.FormatUint(uint64(scheduler.ID), 10)})\">Delete</button></td></tr>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "</tbody></table></div></div><dialog id=\"create-scheduler-modal\" class=\"modal\"><div class=\"modal-box\"><h3 class=\"font-bold text-lg mb-4\">Create Scheduler</h3><form id=\"create-scheduler-form\" onsubmit=\"createScheduler(event)\"><div class=\"form-control\"><label class=\"label\"><span class=\"label-text\">Name</span></label> <input type=\"text\" name=\"name\" class=\"input input-bordered\" required></div><div class=\"form-control\"><label class=\"label\"><span class=\"label-text\">Type</span></label> <select name=\"type\" class=\"select select-bordered\" onchange=\"toggleSchedulerFields(this.value)\" required><option value=\"device\">Device-based</option> <option value=\"date\">Date-based</option></select></div><div class=\"form-control\"><label class=\"label\"><span class=\"label-text\">Device</span></label> <select name=\"device_id\" class=\"select select-bordered\" required>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, device := range props.Devices {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "<option value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var7 string
			templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.FormatUint(uint64(device.ID), 10))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/schedulers.templ`, Line: 101, Col: 65}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 string
			templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(device.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/templates/schedulers.templ`, Line: 101, Col: 79}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "</select></div><div class=\"form-control\" id=\"threshold-field\"><label class=\"label\"><span class=\"label-text\">Threshold</span></label> <input type=\"number\" name=\"threshold\" class=\"input input-bordered\" step=\"0.1\" required></div><div class=\"form-control hidden\" id=\"date-field\"><label class=\"label\"><span class=\"label-text\">Date</span></label> <input type=\"datetime-local\" name=\"date\" class=\"input input-bordered\"></div><div class=\"modal-action\"><button type=\"submit\" class=\"btn btn-primary\">Create</button> <button type=\"button\" class=\"btn\" onclick=\"document.getElementById(&#39;create-scheduler-modal&#39;).close()\">Cancel</button></div></form></div></dialog><script>\n\t\t\t\tfunction toggleSchedulerFields(type) {\n\t\t\t\t\tconst thresholdField = document.getElementById('threshold-field');\n\t\t\t\t\tconst dateField = document.getElementById('date-field');\n\t\t\t\t\tconst thresholdInput = document.querySelector('input[name=\"threshold\"]');\n\t\t\t\t\tconst dateInput = document.querySelector('input[name=\"date\"]');\n\n\t\t\t\t\tif (type === 'device') {\n\t\t\t\t\t\tthresholdField.classList.remove('hidden');\n\t\t\t\t\t\tdateField.classList.add('hidden');\n\t\t\t\t\t\tthresholdInput.required = true;\n\t\t\t\t\t\tdateInput.required = false;\n\t\t\t\t\t} else {\n\t\t\t\t\t\tthresholdField.classList.add('hidden');\n\t\t\t\t\t\tdateField.classList.remove('hidden');\n\t\t\t\t\t\tthresholdInput.required = false;\n\t\t\t\t\t\tdateInput.required = true;\n\t\t\t\t\t}\n\t\t\t\t}\n\n\t\t\t\tasync function createScheduler(event) {\n\t\t\t\t\tevent.preventDefault();\n\t\t\t\t\tconst form = event.target;\n\t\t\t\t\tconst formData = new FormData(form);\n\t\t\t\t\tconst data = {\n\t\t\t\t\t\tname: formData.get('name'),\n\t\t\t\t\t\ttype: formData.get('type'),\n\t\t\t\t\t\tdevice_id: parseInt(formData.get('device_id')),\n\t\t\t\t\t\tthreshold: parseFloat(formData.get('threshold')),\n\t\t\t\t\t\tdate: formData.get('date')\n\t\t\t\t\t};\n\n\t\t\t\t\ttry {\n\t\t\t\t\t\tconst response = await fetch('/api/schedulers', {\n\t\t\t\t\t\t\tmethod: 'POST',\n\t\t\t\t\t\t\theaders: {\n\t\t\t\t\t\t\t\t'Content-Type': 'application/json'\n\t\t\t\t\t\t\t},\n\t\t\t\t\t\t\tbody: JSON.stringify(data)\n\t\t\t\t\t\t});\n\n\t\t\t\t\t\tif (!response.ok) {\n\t\t\t\t\t\t\tthrow new Error('Failed to create scheduler');\n\t\t\t\t\t\t}\n\n\t\t\t\t\t\twindow.location.reload();\n\t\t\t\t\t} catch (error) {\n\t\t\t\t\t\talert('Failed to create scheduler: ' + error.message);\n\t\t\t\t\t}\n\t\t\t\t}\n\n\t\t\t\tasync function deleteScheduler(id) {\n\t\t\t\t\tif (!confirm('Are you sure you want to delete this scheduler?')) {\n\t\t\t\t\t\treturn;\n\t\t\t\t\t}\n\n\t\t\t\t\ttry {\n\t\t\t\t\t\tconst response = await fetch(`/api/schedulers?id=${id}`, {\n\t\t\t\t\t\t\tmethod: 'DELETE'\n\t\t\t\t\t\t});\n\n\t\t\t\t\t\tif (!response.ok) {\n\t\t\t\t\t\t\tthrow new Error('Failed to delete scheduler');\n\t\t\t\t\t\t}\n\n\t\t\t\t\t\twindow.location.reload();\n\t\t\t\t\t} catch (error) {\n\t\t\t\t\t\talert('Failed to delete scheduler: ' + error.message);\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t</script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
