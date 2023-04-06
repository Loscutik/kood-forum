function header() {
  if (document.body.scrollTop == 0 || document.documentElement.scrollTop == 0) {
      document.getElementById("header").style.marginTop = "0px";
  }
  if (document.body.scrollTop > 50 || document.documentElement.scrollTop > 50) {
document.getElementById("header").style.marginTop = "-57px";
  }
}
function header2() {
document.getElementById("header").style.marginTop = "0px";
}
function header3() {
  if (document.body.scrollTop > 50 || document.documentElement.scrollTop > 50) {
document.getElementById("header").style.marginTop = "-57px";
  }
}
function signIn() {
  document.getElementById("darkness").style.display = "block";
  document.getElementById("signinform").style.display = "block";
  document.getElementById("signupform").style.display = "none";
}
function signUp() {
    document.getElementById("warning").style.display = "none";
    document.getElementById("darkness").style.display = "block";
    document.getElementById("signupform").style.display = "block";
    document.getElementById("signinform").style.display = "none";
  
}
function checkPasswords(){
  var password=document.getElementById("password")
  var confirmPassword=document.getElementById("confirm_password")
  confirmPassword.addEventListener("change", (event) => {
    if(password.value!=confirmPassword.value){
      document.getElementById("warning").style.display = "inline";
    }else{
      document.getElementById("warning").style.display = "none";
    }
  });
}
function darkness() {
  document.getElementById("darkness").style.display = "none";
  document.getElementById("signinform").style.display = "none";
  document.getElementById("signupform").style.display = "none";
}
function openFilters() {
  if (document.getElementById("filterform").style.display == "none") {
    document.getElementById("filterform").style.display = "block";
    document.getElementById("filterslogo").style.display = "none";
    document.getElementById("closefilterslogo").style.display = "block";
  } else {
    document.getElementById("filterform").style.display = "none";
    document.getElementById("filterslogo").style.display = "block";
    document.getElementById("closefilterslogo").style.display = "none";
  }
}
function ShowPassword() {
  var x = document.getElementsByClassName("password");
  for (let i = 0; i < x.length; i++) {
    if (x[i].type === "password") {
      x[i].type = "text";
    } else {
      x[i].type = "password";
    }
  }
}