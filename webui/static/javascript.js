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
    document.getElementById("darkness").style.display = "block";
    document.getElementById("signupform").style.display = "block";
    document.getElementById("signinform").style.display = "none";
  
}
function checkForm(){
  const submitButton=document.querySelector("#signup_submit");
  const password=document.getElementById("password");
  const confirmPassword=document.getElementById("confirm_password");
  const warning=document.querySelector("#warning");
  const form=document.querySelector("#signup_form");

  submitButton.addEventListener("click", (event) => {
if (document.querySelector("#email").value==null ||document.querySelector("#username").value==null ||password.value==null ||confirmPassword.value==null|| document.querySelector("#email").value=="null" ||document.querySelector("#username").value=="null" ||password.value=="null" ||confirmPassword.value==""){
      warning.innerHTML="fill all fields";
      warning.style.display="block";
    }else if (password.value!=confirmPassword.value){
      warning.innerHTML="passwords do not match";
      warning.style.display="block";
    }else{
      form.submit();
    }
  });
}

function handleLike(id){
  // needed : "messageType"("posts_likes", "comments_likes") "messageID"(#)  "like"(bool) 
  const clickedElement = document.getElementById(id);
  // parse id
  // form data for post 
  // form option for fetch
  // feth
  // handle the responce - change like numbers

const myHeaders = new Headers();
myHeaders.append("Accept", "image/jpeg");

const myInit = {
  method: "GET",
  headers: myHeaders,
  mode: "cors",
  cache: "default",
};

const myRequest = new Request("flowers.jpg");
fetch(myRequest,myInit)
  .then((response) => {
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    return response.blob();
  })
  .then((response) => {
    myImage.src = URL.createObjectURL(response);
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