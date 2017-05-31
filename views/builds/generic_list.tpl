{{template "header.tpl" .}}

      <link rel="stylesheet" href="{{url "/static/css/theme.bootstrap.css"}}">

      <script type="text/javascript" src="{{url "/static/js/jquery.tablesorter.min.js"}}"></script>
      <script type="text/javascript" src="{{url "/static/js/jquery.tablesorter.pager.min.js"}}"></script>
      <script type="text/javascript" src="{{url "/static/js/jquery.tablesorter.widgets.min.js"}}"></script>

      <script>
        // With customizations
        $(document).ready(function() {
          $.extend($.tablesorter.themes.bootstrap, {
            // these classes are added to the table. To see other table classes available,
            // look here: http://twitter.github.com/bootstrap/base-css.html#tables
            table      : 'table table-bordered',
            caption    : 'caption',
            header     : 'bootstrap-header', // give the header a gradient background
            footerRow  : '',
            footerCells: '',
            icons      : '', // add "icon-white" to make them white; this icon class is added to the <i> in the header
            sortNone   : 'bootstrap-icon-unsorted',
            sortAsc    : 'icon-chevron-up glyphicon glyphicon-chevron-up',     // includes classes for Bootstrap v2 & v3
            sortDesc   : 'icon-chevron-down glyphicon glyphicon-chevron-down', // includes classes for Bootstrap v2 & v3
            active     : '', // applied when column is sorted
            hover      : '', // use custom css here - bootstrap class may not override it
            filterRow  : '', // filter row class
            even       : '', // odd row zebra striping
            odd        : ''  // even row zebra striping
          });

          $("#pkgtable").tablesorter({
            theme: "bootstrap",
            headerTemplate: "{content} {icon}",
            widgets: ["uitheme", "filter", "zebra"],
            textExtraction: {
              4: function(node, table, cellIndex) {
                return $(node).find("div").text();
              },
              5: function(node, table, cellIndex) {
                return $(node).find("img").attr("alt");
              },
              6: function(node, table, cellIndex) {
                  return $(node).attr("original-time");
              }
            },
          });
        });
      </script>

      <div class="page-header">
        <h1>{{.Title}}</h1>
      </div>

      {{if .LoggedIn}}
      <div class="row">
          <form id="bulkForm" class="form-inline">
              <div class="form-group">
                  <p class="form-control-static"><strong>Bulk Actions: </strong></p>
              </div>
              <div class="form-group">
                  <select id="bulkAction" class="form-control" name="action">
                      <option value="Neutral">Clear/Neutral</option>
                      <option value="Down">Reject</option>
                      <option value="Up">Accept</option>
                      <option value="Maintainer">Maintainer Accept (if possible)</option>
                      {{if .QAControls}}
                      <option value=""> --- </option>
                      <option value="QABlock">QA Block</option>
                      <option value="QAPush">QA Push</option>
                      <option value="QAClear">QA Clear</option>
                      {{end}}
                      <option value=""> --- </option>
                      <option value="Finalize">Finalize (if possible)</option>
                  </select>
              </div>
              <div class="form-group">
                  <input type="text" id="bulkComments" class="form-control" placeholder="Comments (Optional)">
              </div>
              <button type="submit" id="bulkApply" class="form-control btn btn-default">Apply</button>
              <div class="form-group pull-right">
                  <a id="bulkSelectAll">Select All</a>
                  <a id="bulkDeselectAll">Deselect All</a>
              </div>
          </form>
      </div>
      {{end}}

      <div class="row table-responsive">
        <table class="table tablesorter" id="pkgtable">
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Submitter</th>
              <th>For</th>
              <th>Karma</th>
              <th>Build Date</th>
            </tr>
          </thead>
          <tbody>
            {{$out := .}}
            {{with .Packages}}
              {{range .}}
              <tr>
                {{if $out.LoggedIn}}
                <td><div class="checkbox" style="margin: 0 0 0 0;"><label><input is-id type="checkbox" name="id" value="{{.Id}}">{{.Id}}</label></td>
                {{else}}
                <td>{{.Id}}</td>
                {{end}}
                <td><a href="{{urldata "/builds/{{.Id}}" .}}">{{.Name}}-{{.SourceEVR}} ({{.Architecture}})</a></td>
                <td>{{.Submitter.Email | emailat}}</td>
                <td>{{.Platform}}/{{.Repo}}</td>
                <td>{{$karma := mapaccess .Id $out.PkgKarma}}<img src="{{if eq $karma "0"}}//img.shields.io/badge/karma-   {{$karma}}-yellow.png{{else}}{{if lt $karma "0"}}//img.shields.io/badge/karma-  -{{$karma}}-orange.png{{else}}{{if gt $karma "0"}}//img.shields.io/badge/karma- +{{$karma}}-yellowgreen.png{{end}}{{end}}{{end}}" alt="{{$karma}}"></td>
                <td data-type="time" original-time="{{.BuildDate | iso8601}}">{{.BuildDate | iso8601}}</td>
              </tr>
              {{end}}
            {{end}}
        </table>
        <center><span class="label label-default">{{.Entries}} {{if eq .Entries 1}}entry{{else}}entries{{end}} returned.</span></center>
      </div>

      <div style="visibility:hidden;display:none;" id="bulkFrames"></div>

      <script type="text/javascript">
      $(function() {
          $("td[data-type='time']").each(function() {
              var node = $(this);
              node.attr("moment-time", moment(node.attr("original-time")).fromNow());
              node.text(node.attr("moment-time"));
              node.hover(function() {
                  node.text(node.attr("original-time"));
              }, function() {
                  node.text(node.attr("moment-time"));
              });
          });

          {{if .LoggedIn}}

          var ids = [];
          var failed = [];
          var current = 0;

          var handleNext = function() {
              if (ids.length === 0) {
                  if (failed.length > 0) {
                      alert("failed to bulk action: " + failed);
                  }

                  window.location.reload(true);
                  return;
              }

              current = ids.shift();
              console.log("running " + current);
              $("#bulkApply").text("Applying " + current + "...");

              var bulkFrame = $("<iframe>", {src: "{{url "/builds/"}}" + current}).on("load", function() {
                  this.contentWindow.$("#voteForm").append($("<input>", {type: 'hidden', name: 'type', value: $("#bulkAction").val()}));
                  this.contentWindow.$("#voteForm").append($("<input>", {type: 'hidden', name: 'comment', value: $("#bulkComments").val()}));
                  this.contentWindow.$("input[type=radio]").remove();
                  this.contentWindow.$("input[type=text]").remove();

                  $(this).off("load");
                  $(this).on("load", function() {
                      handleNext();
                      bulkFrame.remove();
                  });
                  this.contentWindow.$("#voteForm").submit();
              });

              $("#bulkFrames").append(bulkFrame);
          };

          var running = false;

          // for each checkbox, post to builds with the desired outcome
          $("#bulkForm").submit(function(event) {
              event.preventDefault();

              if (running) { // we don't want to allow it to run again
                  return;
              }

              if ($("#bulkAction").val() === '') {
                  return;
              }

              $("[is-id]:checked").each(function() {
                  ids.push($(this).val());
              });

              $("#bulkApply").prop('disabled', true);
              $("#bulkApply").text("Please wait...");

              running = true;

              handleNext();
          });

          $("#bulkSelectAll").click(function() {
              $("[is-id]").each(function() {
                  $(this).prop('checked', true);
              });
          });

          $("#bulkDeselectAll").click(function() {
              $("[is-id]").each(function() {
                  $(this).prop('checked', false);
              });
          });

          {{end}}
      });
      </script>

{{template "footer.tpl" .}}
