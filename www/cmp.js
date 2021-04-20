// Get the expected elements based on common names.
var saltElement = document.getElementById('salt-container');
var saltEmailField = document.getElementById('email');
var saltUpdateButton = document.getElementById('update');
var saltResetButton = document.getElementById('reset-salt');
var saltField = document.getElementById('salt');

// Record the original message in the update button.
var saltInnerHTML = saltUpdateButton.innerHTML;

// Initialize the salt element and user interface.
var saltUI = new SWANSalt(saltElement, saltField.value);

// Change the update button to explain 4 icons need to be selected and disable 
// it.
function saltReset() {
    saltUpdateButton.innerHTML = 'Choose 4 Icons';
    saltUpdateButton.setAttribute('disabled', '');
}

// Show the salt form group.
function saltShow() {
    $('#salt-form-group').collapse('show');
}

// Clear the salt 
function saltClear() {
    $('#salt-form-group').collapse('hide');
    complete();
}

// Called when the salt grid has been completed.
function saltComplete() {
    saltField.value = saltUI.stringValue;
    saltUpdateButton.innerHTML = saltInnerHTML;
    saltUpdateButton.removeAttribute('disabled'); 
}

function saltEmailChange(e) {
    var email = e.value;
    if(email == '') {
        // If there is no email address then clear the salt.
        saltClear();
    } else if (/^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/.test(email)) {
        // If there is a valid email then enable the 
        saltShow();
        saltReset();
    }
}

// If the email changes then trigger the display or removal of the salt element.
saltEmailField.addEventListener('input', (event) => {
    saltEmailChange(event.target);
});

// If the reset button is clicked then reset the salt and the update button.
saltResetButton.addEventListener('click', (event) => {
    saltUI.reset();
    saltReset();
});

// When the salt UI completes call this function.
saltUI.onComplete(saltComplete);

// Trigger the display of the salt grid if the email address is present.
saltEmailChange(saltEmailField);

// Set the salt to complete if the salt is already provided.
if (saltField.value != "") {
    saltComplete();
}