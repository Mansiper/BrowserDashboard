window.onload = function() {
  $('#show_all_messages').on('change', function() {
    onAllMessagesClick();
    saveOptions();
  });
  $('#show_main_messages').on('change', function() { saveOptions(); });
  $('#show_test_messages').on('change', function() { saveOptions(); });
  $('#show_asev_messages').on('change', function() { saveOptions(); });
  $('#show_db_messages').on('change', 	function() { saveOptions(); });

  loadOptions();
  onAllMessagesClick();
};

function onAllMessagesClick() {
  var checked = $('#show_all_messages')[0].checked;
  $('#show_main_messages')[0].disabled = checked ? '' : 'disabled';
  $('#show_test_messages')[0].disabled = checked ? '' : 'disabled';
  $('#show_asev_messages')[0].disabled = checked ? '' : 'disabled';
  $('#show_db_messages')[0].disabled = 	 checked ? '' : 'disabled';
}

function saveOptions() {
  window.localStorage.setItem('show_all_messages',  $('#show_all_messages')[0].checked);
  window.localStorage.setItem('show_main_messages', $('#show_main_messages')[0].checked);
  window.localStorage.setItem('show_test_messages', $('#show_test_messages')[0].checked);
  window.localStorage.setItem('show_asev_messages', $('#show_asev_messages')[0].checked);
  window.localStorage.setItem('show_db_messages',   $('#show_db_messages')[0].checked);
}

function loadOptions() {
  $('#show_main_messages')[0].checked = window.localStorage.getItem('show_main_messages') == "true";
  $('#show_test_messages')[0].checked = window.localStorage.getItem('show_test_messages') == "true";
  $('#show_asev_messages')[0].checked = window.localStorage.getItem('show_asev_messages') == "true";
  $('#show_db_messages')[0].checked   = window.localStorage.getItem('show_db_messages')   == "true";
  $('#show_all_messages')[0].checked  = window.localStorage.getItem('show_all_messages')  == "true";
}