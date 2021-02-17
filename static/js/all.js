var button = $('.button'),
		spinner = '<span class="spinner"></span>';

button.click(function() {
	if (!button.hasClass('loading')) {
		button.toggleClass('loading').html(spinner);
	}
	else {
		button.toggleClass('loading').html("Load");
	}
})

/*Dropdown Menu*/
$('.dropdown').click(function () {
        $(this).attr('tabindex', 1).focus();
        $(this).toggleClass('active');
        $(this).find('.dropdown-menu').slideToggle(300);
    });
    $('.dropdown').focusout(function () {
        $(this).removeClass('active');
        $(this).find('.dropdown-menu').slideUp(300);
    });
    $('.dropdown .dropdown-menu li').click(function () {
        $(this).parents('.dropdown').find('span').text($(this).text());
        $(this).parents('.dropdown').find('input').attr('value', $(this).attr('id'));
    });
/*End Dropdown Menu*/


$('.dropdown-menu li').click(function () {
  var input = '<strong>' + $(this).parents('.dropdown').find('input').val() + '</strong>',
      msg = '<span class="msg">Hidden input value: ';
  $('.msg').html(msg + input + '</span>');
});

//$('.loader').hide();
$('#text').show();
$('#loaded').hide();

function myFunction() {
	$('#text').hide();
	//$('.loader').show();
	$('#loaded').show();
}
