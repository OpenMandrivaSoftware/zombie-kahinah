$("document").ready(function(){

  $("#logout").hide();

  function loggedIn(email){
    $("#login").hide();
    $("#logout").show();
    $("#persona-user").text(email);
  }

  function loggedOut(){
    $("#logout").hide();
    $("#login").show();
    $("#persona-user").text("");
  }

  var user = null;

  $.ajax({
    url: window.urlPrefix + '/auth/check',
    success: function(res, status, xhr) {
      if (res === "") {
        loggedOut();
      }
      else {
        loggedIn(res);
        user = res;
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
