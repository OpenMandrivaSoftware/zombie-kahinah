{{template "header.tpl" .}}

      <!-- Jumbotron -->
      <div class="jumbotron">
        <h1>Kahinah, the OpenMandriva Update System</h1>
        <p class="lead">This is Kahinah v3 (now termed Zombie-Kahinah), which hooks into <a href="//abf.openmandriva.org">ABF</a>, allowing developers and daredevil users to test updates before they are pushed to live.</p>
        <p><a class="btn btn-lg btn-success" href="{{url "/builds/testing"}}" role="button">Recent Builds</a></p>
      </div>

      <!-- Infos -->
      <div class="row">
        <div class="col-lg-6">
          <h2>This is an unsupported tool.</h2>
          <p class="text-danger">While this tool has been extensively used, it is scheduled to be deprecated. Caution is advised. If any updates are not pushed or go missing, please alert OpenMandriva QA.</p>
          <p><a class="btn btn-primary" href="mailto:qa@openmandriva.org" role="button">Contact &raquo;</a></p>
        </div>
        <div class="col-lg-6">
          <h2>News</h2>
          <p>Not available.</p>
       </div>
      </div>

{{template "footer.tpl" .}}
