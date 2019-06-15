var asevTime = new Date(2000, 0, 1, 0, 0);
var connDbTime = new Date(2000, 0, 1, 0, 0);
var errorTime = new Date(2000, 0, 1, 0, 0);
var oldJobj = { //Значения по умолчанию - всё хорошо.
  'back_check_1':1,
  'back_check_2':1,
  'backtest_check_1':1,
  'backtest_check_2':1,
  'api_check':1,
  'apitest_check':1,
  'api2_check':1,
  'api2test_check':1,
  'gps_check':1,
  'as_check':1,
  'asev_check':1,
  'dbmain_check':1,
  'dbtest_check':1
};

var opts = {
  show_all_msg:  window.localStorage.getItem('show_all_messages')  == 'true',
  show_main_msg: window.localStorage.getItem('show_main_messages') == 'true',
  show_test_msg: window.localStorage.getItem('show_test_messages') == 'true',
  show_asev_msg: window.localStorage.getItem('show_asev_messages') == 'true',
  show_db_msg: 	 window.localStorage.getItem('show_db_messages')   == 'true',
};

function loadData() {
	$.ajax({
		type: 'GET',
		contentType: 'application/json; charset=UTF-8',
		url: 'https://companyurl:123/dbdata',
		global: true,
		success: function(data, textStatus, jqXHR) {
			var jobj = JSON.parse(data);
			if (jobj == null) return;
			$('#connError').addClass('hidden');

			var val, elems, elem, serviceName;
			var allGood = 0;
			for (var key in jobj) {
				val = jobj[key];
				if (key.indexOf('DB_', 0) != 0 && val == 0)
					allGood += 1;
				elems = $('#'+key);
				if (elems.length != 0)
					elem = elems[0];
				if (key.indexOf('DB_', 0) == 0) {
					if (elem) {
						if (key == 'DB_Threads_running')
							$(elem).text(val + ' /');
						else $(elem).text(val);
					}
					
					var serviceMsg = '';
					if (key == 'DB_Threads_running' && val > 30 && opts['show_db_msg'])
						serviceMsg = 'Слишком много активных подключений к БД.';
					else
						serviceMsg = '';
					
					if (opts['show_all_msg'] && serviceMsg && serviceMsg != '' &&
					((Date.now() - connDbTime > 60000/*1 минута*/ && key == 'DB_Threads_running') || key != 'DB_Threads_running')) {
						chrome.notifications.clear(key, function() {});
						chrome.notifications.create(key, {
							iconUrl: chrome.runtime.getURL('icon_bad.png'),
							title: 'Проблема!',
							type: 'basic',
							message: serviceMsg,
							priority: 2,
						}, function() {});
						if (key == 'DB_Threads_running')
							connDbTime = Date.now();
					}					
				}
				else {
					if (val == 1) {
						if (elem) {
							$(elem).removeClass('red');
							$(elem).removeClass('yellow');
							$(elem).addClass('green');
						}
						
						chrome.notifications.clear(key, function() {});
					} else {
						if (elem) {
							$(elem).removeClass('green');
							$(elem).removeClass('yellow');
							$(elem).addClass('red');
						}
						
						var serviceName = '';
						switch (key) {
							case 'back_check_1':
              case 'back_check_2':
								serviceName = opts['show_main_msg'] ? 'Бэком' : null;
								break;
							case 'api_check':
								serviceName = opts['show_main_msg'] ? 'API' : null;
								break;
							case 'api2_check':
								serviceName = opts['show_main_msg'] ? 'API2' : null;
								break;
							case 'gps_check':
								serviceName = opts['show_main_msg'] ? 'сервисом GPS' : null;
								break;
							case 'as_check':
								serviceName = opts['show_main_msg'] ? 'сервисом звонков' : null;
								break;
							case 'asev_check':
								serviceName = opts['show_asev_msg'] ? 'сервисом событий звонков (либо просто давно никто не звонил)' : null;
								break;
							case 'dbmain_check':
								serviceName = opts['show_main_msg'] ? 'боевой базой данных' : null;
								break;
              case 'backtest_check_1':
              case 'backtest_check_2':
								serviceName = opts['show_test_msg'] ? 'тестовым Бэком' : null;
								break;
              case 'apitest_check':
								serviceName = opts['show_test_msg'] ? 'API test' : null;
								break;
              case 'api2test_check':
								serviceName = opts['show_test_msg'] ? 'API2 test' : null;
								break;
              case 'dbtest_check':
                serviceName = opts['show_test_msg'] ? 'тестовой базой данных' : null;
                break;
							default:
								serviceName = '';
								break;
						}
						
						if (opts['show_all_msg'] && serviceName && serviceName != '' && val != oldJobj[key] &&
						((Date.now() - asevTime > 60000/*1 минута*/ && key == 'asev_check') || key != 'asev_check')) {
							chrome.notifications.clear(key, function() {});
							chrome.notifications.create(key, {
								iconUrl: chrome.runtime.getURL('icon_bad.png'),
								title: 'Потеря связи!',
								type: 'basic',
								message: 'Пропала связь с ' + serviceName + '!',
								priority: 2,
							}, function() {});
							if (key == 'asev_check')
								asevTime = Date.now();
						}
					}
				}
			}
			if (allGood === 0)
				chrome.browserAction.setIcon({path:'icon_good.png'});
			else chrome.browserAction.setIcon({path:'icon_bad.png'});
			chrome.notifications.clear('CompanyServiceDisconnected', function() {});

      oldJobj = jobj;
		},
		error: function(jqXHR, textStatus, errorThrown) {
			elem = $('#connError')[0];
			$(elem).removeClass('hidden');
			chrome.browserAction.setIcon({path:'icon_def.png'});
			
			if (Date.now() - errorTime > 300000/*5 минут*/) {
				chrome.notifications.create('CompanyServiceDisconnected', {
					iconUrl: chrome.runtime.getURL('icon_def.png'),
					title: 'Потеря связи!',
					type: 'basic',
					message: 'Пропала связь с сервисом статистики. Обратитесь к разработчикам для её восстановления.',
					priority: 2,
				}, function() {});
				errorTime = Date.now();
			}
		}
	});
}

window.onload = function() {
	loadData();
	setInterval(loadData, 2000);
}