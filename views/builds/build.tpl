{{template "header.tpl" .}}

      <div class="row">
          <br/>
          {{if eq .Package.Status "testing"}}<div class="card bg-light border-warning">{{else}}
          {{if eq .Package.Status "rejected"}}<div class="card bg-light border-danger">{{else}}
          {{if eq .Package.Status "published"}}<div class="card bg-light border-success">{{else}}
          <div class="card bg-light border-primary">{{end}}{{end}}{{end}}
            <h2 class="card-header d-flex justify-content-between">
              {{.Package.Name}} <small>{{.SourceEVR}} ({{.Package.Architecture}}) {{.Package.BuildDate.Year}}-{{.Package.Id}}</small>
              <div class="text-end">{{.Karma}} {{if .KarmaControls}}<a href="#" class="btn" data-bs-toggle="modal" data-bs-target="#voteModal"><i class="fa fa-2x {{if .UserVote}}fa-check-square-o{{else}}fa-pencil-square-o{{end}}"></i></a>{{end}}</div>
            </h2>
            <table class="table table-condensed">
              <tbody>
                <tr>
                  <td><b>Status</b></td>
                  <td>{{.Package.Status}}</td>
                </tr>
                <tr>
                  <td><b>Submitter</b></td>
                  <td>{{.Package.Submitter.Email | emailat}}</td>
                </tr>
                <tr>
                  <td><b>Platform<b></td>
                  <td>{{.Package.Platform}}</td>
                </tr>
                <tr>
                  <td><b>Repository<b></td>
                  <td>{{.Package.Repo}}</td>
                </tr>
                <tr>
                  <td><b>URL<b></td>
                  <td><a href="{{.Url}}">{{.Url}}</a></td>
                </tr>
                <tr>
                  <td><b>Packages<b></td>
                  <td>
                    <table class="table table-bordered table-condensed table-responsive">
                      <tbody>
                        {{with .Package.Packages}}
                          {{range .}}
                        <tr>
                          <td>{{.Name}}-{{if gt .Epoch 0}}{{.Epoch}}:{{end}}{{.Version}}-{{.Release}}{{if .Arch}}.{{.Arch}}{{end}}.{{.Type}}</td>
                        </tr>
                          {{end}}
                        {{end}}
                      </tbody>
                    </table>
                  </td>
                </tr>
                <tr>
                  <td><b>Build Date</b></td>
                  <td>{{.Package.BuildDate}}</td>
                </tr>
                <tr>
                  <td><b>Last Updated</b></td>
                  <td>{{.Package.Updated}}</td>
                </tr>
            </table>
          </div>
      </div>

      <!-- diff & changelog -->
      <div class="row">

          <div class="card border-success">
            <h5 class="card-header"><button class="btn btn-outline-success" data-bs-toggle="collapse" href="#diff">Git Diff</button></h5>
            <div id="diff" class="card-collapse collapse">
              <pre class="brush: diff">{{.Package.Diff}}</pre>
            </div>
          </div>

          <div class="card border-info">
            <h5 class="card-header"><button class="btn btn-outline-info" data-bs-toggle="collapse" href="#cnlog">Changelog</button></h5>
            <div id="cnlog" class="card-collapse collapse">
              <pre class="pre-scrollable">{{if .Changelog}}{{.Changelog}}{{else}}Not Available{{end}}</pre>
            </div>
          </div>

          <div class="card card-default">
            <h4 class="card-header"><button class="btn btn-default">Votes</button></h4>
            <div class="card-body">
              {{if .Votes}}
              <table class="table table-condensed table-responsive table-bordered">
                {{with .Votes}}
                  {{range .}}
                <tr class="{{if eq .Value 1}}bg-success text-white{{end}}{{if eq .Value 2}}bg-danger text-white{{end}}{{if eq .Value 3}}bg-info{{end}}{{if eq .Value 4}}bg-warning{{end}}"><td>{{.Key.User.Email | emailat}}</td><td>{{if .Key.Comment}}{{.Key.Comment}}{{else}}<em>No Comment.</em>{{end}}</td><td>{{if .Key.Time}}{{.Key.Time | since}}{{else}}[voted before timekeeping began]{{end}}</td></tr>
                  {{end}}
                {{end}}
              </table>
              {{else}}
              Nobody has submitted a vote for this package yet. Go test it, and vote!
              {{end}}
            </div>
          </div>

      </div>

      <!-- Vote Modal -->

      <div class="modal fade" id="voteModal" tabindex="-1" role="dialog">
        <div class="modal-dialog modal-dialog-centered">
          <center>
            <div class="modal-content">
              <form class="form-inline" role="form" method="post" id="voteForm">
                {{ .xsrf_data }}
                <div class="modal-header d-flex justify-content-between">
                  <button type="button" class="btn-close m-0" data-bs-dismiss="modal" aria-label="Close"></button>
                  <h4 class="modal-title" id="Vote Modal">Submit Vote</h4>
                </div>
                <div class="modal-body">
                  <div class="btn-group" data-toggle="buttons">                  
                    <input type="radio" class="btn-check" name="type" value="Neutral" id="opNeutral" {{if eq .UserVote 0}}checked{{end}}>
                    <label class="btn-outline-primary" for="opNeutral"><i class="fa fa-lg fa-meh-o"></i><br>No Vote</label>

                    <input type="radio" class="btn-check" name="type" value="Down" id="opDown" {{if eq .UserVote -1}}checked{{end}}>
                    <label class="btn btn-outline-danger" for="opDown"><i class="fa fa-lg fa-frown-o"></i> Reject</label>

                    <input type="radio" class="btn-check" name="type" value="Up" id="opUp" {{if eq .UserVote 1}}checked{{end}}>
                    <label class="btn btn-outline-success" for="opUp"><i class="fa fa-lg fa-smile-o"></i> Accept</label>

                    {{if .MaintainerControls}}
                    {{if .MaintainerTime}}<input type="radio" class="btn-check" name="type" value="Maintainer" id="opMaintainer" {{if eq .UserVote 2}}checked{{end}}>{{end}}
                    <label class="btn btn-outline-primary" for="opMaintainer" {{if not .MaintainerTime}}disabled="disabled"{{end}}><i class="fa fa-lg fa-thumbs-o-up"></i> Maintainer Push</label>
                    {{end}}

                    {{if .QAControls}}
                    <input type="radio" class="btn-check" name="type" id="voteQADown" value="QABlock">
                    <label class="btn btn-outline-warning" for="voteQADown"><i class="fa fa-lg fa-thumbs-o-down"></i> QA Block</label>

                    <input type="radio" class="btn-check" name="type" id="voteQAUp" value="QAPush">
                    <label class="btn btn-outline-warning" for="voteQAUp"><i class="fa fa-lg fa-thumbs-o-up"></i> QA Push</label>

                    <input type="radio" class="btn-check" name="type" id="voteQAClear" value="QAClear">
                    <label class="btn btn-outline-warning" for="voteQAClear"><i class="fa fa-lg fa-scissors"></i> QA Clear</label>
                    {{end}}
                    {{if .FinalizeControls}}
                    <input type="radio" class="btn-check" name="type" id="voteFinalize" value="Finalize">
                    <label class="btn btn-outline-info" for="voteFinalize" ><i class="fa fa-lg fa-flag-o"></i> Finalize</label>
                    {{end}}
                  </div>
                </div>
                {{if .MaintainerControls}}{{if not .MaintainerTime}}
                <div class="alert alert-info"><b>This is your update!</b> Unfortunately, you need to wait {{.MaintainerHoursNeeded}} hours since the Build Date until you can activate Maintainer Accept.</div>
                {{else}}<div class="alert alert-info"><b>This is your update!</b> You can activate Maintainer Accept now.</div>{{end}}{{end}}
                <div id="voteModalAlertPlaceholder"></div>
                <div class="modal-body">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-lg fa-comment-o"></i> Comment</span>
                    <input type="text" class="form-control" name="comment" placeholder="It's recommended to say something." {{if .KarmaCommentPrev}}value="{{.KarmaCommentPrev}}"{{end}}>
                  </div>
                </div>
                <div class="modal-footer">
                  <button type="button" class="btn btn-default" data-bs-dismiss="modal">Close</button>
                  <button type="submit" type="button" class="btn btn-primary">Submit</button>
                </div>
              </form>
            </div>
          </center>
        </div>
      </div>

      <script>
        $("input").change(function() {
          if ($("#voteQADown").is(':checked')) {
            $('#voteModalAlertPlaceholder').html('<div class="alert alert-danger"><b>Heads up!</b> This adds -9999 karma and allows the build to be finalized immediately!</div>');
          } else if ($("#voteQAUp").is(':checked')) {
            $('#voteModalAlertPlaceholder').html('<div class="alert alert-warning"><b>Heads up!</b> This adds 9999 karma and allows the build to be finalized immediately!</div>');
          } else if ($("#voteQAClear").is(':checked')) {
            $('#voteModalAlertPlaceholder').html('<div class="alert alert-warning"><b>Beware:</b> This removes this update from being controlled by Kahinah. Useful if someone has bypassed via ABF.</div>');
          } else if ($("#voteFinalize").is(':checked')) {
            $('#voteModalAlertPlaceholder').html('<div class="alert alert-warning">After you finalize this update, it will be pushed to the main updates repositories. No more changes can be made.</div>');
          } else {
            $('#voteModalAlertPlaceholder').html('')
          }
        }).change();
      </script>

      <link href="{{url "/static/css/shCore.css"}}" rel="stylesheet" type="text/css" />
      <link href="{{url "/static/css/shThemeDefault.css"}}" rel="stylesheet" type="text/css" />

      <script src="{{url "/static/js/shCore.js"}}" type="text/javascript"></script>
      <script src="{{url "/static/js/shBrushDiff.js"}}" type="text/javascript"></script>
      <script>SyntaxHighlighter.all();</script>

{{template "footer.tpl" .}}
