$("document").ready(function(){

  $("#logout").hide();

  function loggedIn(){
    $("#login").hide();
    $("#logout").show();
  }

  function loggedOut(){
    $("#logout").hide();
    $("#login").show();
  }

  $.ajax({
    url: window.urlPrefix + '/auth/check',
    success: function(res, status, xhr) {
      if (res === "") {
        loggedOut();
      }
      else {
        loggedIn();
      }
    },
    async: false
  });

  $("#login").on("click", function(e) {
    e.preventDefault();
    $("#login").text("Logging in...");
    $("#login").attr("class", "btn btn-info navbar-btn");
    window.location.href = window.urlPrefix + '/auth/login';
  });

  $("#logout").on("click", function(e) {
    e.preventDefault();
    $.get(window.urlPrefix + '/auth/logout', function() {
      window.location.href = window.location.href;
    });
  });

});
