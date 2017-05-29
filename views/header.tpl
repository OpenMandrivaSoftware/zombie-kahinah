<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="">
    <meta name="author" content="">

    <meta name="_xsrf" content="{{.xsrf_token}}" />

    <link rel="shortcut icon" href="{{url "/static/img/favicon.png"}}">

    <title>{{.Title}} | Kahinah v3 on Jasper</title>

    <!-- Bootstrap core CSS -->
    <link href="//netdna.bootstrapcdn.com/bootswatch/3.1.1/spacelab/bootstrap.min.css" rel="stylesheet">

    <!-- Font Awesome -->
    <link href="//netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.min.css" rel="stylesheet">

    <link href="{{url "/static/css/justified-nav.css"}}" rel="stylesheet">

    <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
      <script src="//oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
      <script src="//oss.maxcdn.com/libs/respond.js/1.3.0/respond.min.js"></script>
    <![endif]-->

    <!-- Bootstrap core JavaScript -->
    <script>window.urlPrefix = "{{url ""}}";</script>
    <script src="//code.jquery.com/jquery-2.0.3.min.js"></script>
    <script src="{{url "/static/js/xsrf.js"}}"></script>
    <script src="{{url "/static/js/persona.js"}}"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/moment.js/2.18.1/moment-with-locales.min.js"></script>

  </head>

  <body>

    <div class="container">

      <nav class="navbar navbar-default" role="navigation">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#navbar-collapse">
            <span class="sr-only">Toggle Navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a href="{{url "/"}}" class="navbar-brand"><i class="fa fa-tasks"></i> Kahinah</a>
        </div>

        <div class="collapse navbar-collapse" id="navbar-collapse">
          <ul class="nav navbar-nav">
            <li{{if eq .Loc 0}} class="active" {{end}}><a href="{{url "/"}}">Home</a></li>

            <li {{if eq .Loc 1}}class="active" {{end}}class="dropdown"> <!-- builds -->
              <a href="#" class="dropdown-toggle" data-toggle="dropdown">Builds <b class="caret"></b></a>
              <ul class="dropdown-menu">
                <li{{if eq .Loc 1}}{{if eq .Tab 1}} class="active"{{end}}{{end}}><a href="{{url "/builds/testing"}}">Testing</a></li>
                <li{{if eq .Loc 1}}{{if eq .Tab 2}} class="active"{{end}}{{end}}><a href="{{url "/builds/published"}}">Published</a></li>
                <li{{if eq .Loc 1}}{{if eq .Tab 3}} class="active"{{end}}{{end}}><a href="{{url "/builds/rejected"}}">Rejected</a></li>
                <li{{if eq .Loc 1}}{{if eq .Tab 4}} class="active"{{end}}{{end}}><a href="{{url "/builds"}}">All</a></li>
              </ul>
            </li>

          </ul>

          <!-- login -->
          <div class="navbar-right">
            <a href="{{url "/user"}}"><p class="navbar-text" id="persona-user"></p></a>
            {{if .xsrf_token}}
            <button type="button" class="btn btn-primary navbar-btn" style="display: none" id="login">Github Login</button><button type="button" class="btn btn-warning navbar-btn" style="display: none" id="logout">Logout</button>
            {{end}}
          </div>
        </div>
      </nav>

      {{if .flash.error}}<div class="alert alert-danger">{{.flash.error}}</div>{{end}}
      {{if .flash.warning}}<div class="alert alert-warning">{{.flash.warning}}</div>{{end}}
      {{if .flash.notice}}<div class="alert alert-success">{{.flash.notice}}</div>{{end}}
