{{template "header.tpl" .}}

      <div class="page-header">
        <h1>{{.Title}}</h1>
      </div>

      <div class="row table-responsive">
        <table class="table">
          <thead>
            <tr>
              <th>Update ID</th>
              <th>Name</th>
              <th>Submitter</th>
              <th>For</th>
              <th>Status</th>
              <th>Updated</th>
            </tr>
          </thead>
          <tbody>
            {{with .Packages}}
              {{range .}}
              <tr>
                <td><a href="{{urldata "/builds/{{.Id}}" .}}">{{.BuildDate.Year}}-{{.Id}}</a></td>
                <td><a href="{{urldata "/builds/{{.Id}}" .}}">{{.Name}}-{{.SourceEVR}} ({{.Architecture}})</a></td>
                <td>{{.Submitter.Email | emailat}}</td>
                <td>{{.Platform}}/{{.Repo}}</td>
                <td><img src="{{if eq .Status "testing"}}//img.shields.io/badge/status-TESTING-yellow.png{{else}}
                    {{if eq .Status "rejected"}}//img.shields.io/badge/status-REJECTED-red.png{{else}}
                    {{if eq .Status "published"}}//img.shields.io/badge/status-PUBLISHED-brightgreen.png{{else}}
                    //img.shields.io/badge/status-UNKNOWN-lightgrey.png{{end}}{{end}}{{end}}" alt="{{.Status}}"></td>
                <td data-type="time" original-time="{{.Updated | iso8601}}">{{.Updated | iso8601}}</td>
              </tr>
              {{end}}
            {{end}}
        </table>
      </div>
      <div class="row">
        <div class="col-md-4 col-md-offset-4">
          <form name="input" method="get">
            <div class="input-group">
              <span class="input-group-btn">
                <a href="?page={{.PrevPage}}"><button class="btn btn-default" type="button">&larr;</button></a>
              </span>
              <span class="input-group-addon">Page</span>
              <input type="text" name="page" class="form-control" placeholder="{{.Page}}">
              <span class="input-group-addon">/ {{.Pages}}</span>
              <span class="input-group-btn">
                <a href="?page={{.NextPage}}"><button class="btn btn-default" type="button">&rarr;</button></a>
              </span>
            </div>
          </form>
        </div>
      </div>

      <script type="text/javascript">
        $("td[data-type='time']").each(function() {
            var node = $(this);
            node.attr("original-time", node.text());
            node.attr("moment-time", moment(node.text()).fromNow());
            node.text(node.attr("moment-time"));
            node.hover(function() {
                node.text(node.attr("original-time"));
            }, function() {
                node.text(node.attr("moment-time"));
            });
        });
      </script>

{{template "footer.tpl" .}}
