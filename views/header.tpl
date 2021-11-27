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

    <title>{{.Title}} | Kahinah</title>

    <!-- Bootstrap core CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">

    <!-- Font Awesome -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@fortawesome/fontawesome-free@5.15.4/css/all.min.css" integrity="sha256-mUZM63G8m73Mcidfrv5E+Y61y7a12O5mW4ezU3bxqW4=" crossorigin="anonymous">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@fortawesome/fontawesome-free@5.15.4/css/v4-shims.min.css" integrity="sha256-j+Lxy3vEHGQK0+okRJz6G6UpHhbbu6sO9hv+Q/MhKRA=" crossorigin="anonymous">

    <link href="{{url "/static/css/justified-nav.css"}}" rel="stylesheet">

    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/jquery-ui-dist@1.12.1/jquery-ui.min.css" integrity="sha256-rByPlHULObEjJ6XQxW/flG2r+22R5dKiAoef+aXWfik=" crossorigin="anonymous">

    <!-- Bootstrap core JavaScript -->
    <script>window.urlPrefix = "{{url ""}}";</script>
    <script src="https://cdn.jsdelivr.net/npm/jquery@3.6.0/dist/jquery.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
    <script src="{{url "/static/js/xsrf.js"}}"></script>
    <script src="https://cdn.jsdelivr.net/npm/moment@2.29.1/min/moment-with-locales.min.js" integrity="sha256-E3Snwx6F4t7DiA/L3DgPk6In2M1747JSau+3PWjtS5I=" crossorigin="anonymous"></script>
  </head>

  <body>

  <div class="col-lg-10 mx-auto p-3 py-md-5">
    <nav class="navbar navbar-expand-lg navbar-light bg-light rounded" aria-label="Kahinah Navbar">
      <div class="container-fluid">
        <a class="navbar-brand" href="{{url "/"}}"><i class="fas fa-server"></i> Kahinah</a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbars" aria-controls="navbars" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="navbars">
          <ul class="navbar-nav me-auto mb-2 mb-lg-0">
            <li class="nav-item">
              <a class="nav-link {{if eq .Loc 0}}active{{end}}" href="{{url "/"}}">Home</a>
            </li>
            <li class="nav-item dropdown">
              <a class="nav-link dropdown-toggle {{if eq .Loc 1}}active{{end}}" href="#" id="dropdown" data-bs-toggle="dropdown" aria-expanded="false">Builds</a>
              <ul class="dropdown-menu" aria-labelledby="dropdown">
                <li><a class="dropdown-item {{if eq .Loc 1}}{{if eq .Tab 1}}active{{end}}{{end}}" href="{{url "/builds/testing"}}">Testing</a></li>
                <li><a class="dropdown-item {{if eq .Loc 1}}{{if eq .Tab 2}}active{{end}}{{end}}" href="{{url "/builds/published"}}">Published</a></li>
                <li><a class="dropdown-item {{if eq .Loc 1}}{{if eq .Tab 3}}active{{end}}{{end}}" href="{{url "/builds/rejected"}}">Rejected</a></li>
                <li><a class="dropdown-item {{if eq .Loc 1}}{{if eq .Tab 4}}active{{end}}{{end}}" href="{{url "/builds"}}">All</a></li>
              </ul>
            </li>
            <li class="nav-item">
              <a class="nav-link {{if eq .Loc 2}}active{{end}}" href="{{url "/audit"}}">Audit Log</a>
            </li>
          </ul>

          <div class="d-flex pe-2">
            <span class="navbar-text">{{.user_login}}</span>
          </div>
          <div class="d-flex">
            {{if .xsrf_token}}
              {{if .LoggedIn}}
                <a class="btn btn-sm btn-outline-secondary" href="{{url "/auth/login"}}" id="logout">Logout</a>
              {{else}}
                <a class="btn btn-sm btn-outline-primary" href="{{url "/auth/logout"}}" id="login">Login with Github</a>
              {{end}}
            {{end}}
          </div>
        </div>
      </div>
    </nav>
    <main class="pt-5">

      {{if .flash.error}}<div class="alert alert-danger" role="alert">{{.flash.error}}</div>{{end}}
      {{if .flash.warning}}<div class="alert alert-warning" role="alert">{{.flash.warning}}</div>{{end}}
      {{if .flash.notice}}<div class="alert alert-success" role="alert">{{.flash.notice}}</div>{{end}}
