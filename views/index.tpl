{{template "header.tpl" .}}

      <!-- Jumbotron -->
      <div class="jumbotron">
        <h1>Kahinah, the OpenMandriva Update Tracking System</h1>
        <p class="lead">This is Kahinah v3 (now termed Zombie-Kahinah), tracking builds from <a href="https://abf.openmandriva.org">ABF</a>.</p>
        <p><a class="btn btn-lg btn-success" href="{{url "/builds/testing"}}" role="button">Recent Builds</a></p>
      </div>

      <!-- Infos -->
      <div class="row">
        <div class="col-lg-6">
          <h2>This tool is lightly maintained.</h2>
          <p>You may experience issues using this tool; please report any issues you find to OpenMandriva's QA team.</p>
        </div>
        <div class="col-lg-6">
          <h2>News</h2>
          <p>2021-12: This tool underwent some updates. Apologies for any broken styling you might see.</p>
       </div>
      </div>

{{template "footer.tpl" .}}
