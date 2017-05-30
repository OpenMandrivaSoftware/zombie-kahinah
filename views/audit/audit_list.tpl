{{template "header.tpl" .}}

      <div class="page-header">
        <h1>{{.Title}}</h1>
      </div>

      <div class="row table-responsive">
        <table class="table">
          <thead>
            <tr>
              <th>Update</th>
              <th>User</th>
              <th>Vote</th>
              <th>Comment</th>
              <th>Committed</th>
            </tr>
          </thead>
          <tbody>
            {{with .Karma}}
              {{range .}}
              <tr>
                <td><a href="{{urldata "/builds/{{.List.Id}}" .}}">{{.List.BuildDate.Year}}-{{.List.Id}}: {{.List.Name}}/{{.List.Architecture}}</a></td>
                <td>{{.User.Email | emailat}}</td>
                <td>{{.Vote | convertKarma}}</td>
                <td>{{.Comment}}</td>
                <td data-type="time" original-time="{{.Time | iso8601}}">{{.Time | iso8601}}</td>
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
