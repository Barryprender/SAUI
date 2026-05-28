var amountWrap = document.querySelector('.amount-wrap');

document.querySelectorAll('.btn-preset').forEach(function (btn) {
	btn.addEventListener('click', function () {
		document.getElementById('amount_euros').value = btn.dataset.amount;
		if (amountWrap) {
			amountWrap.classList.remove('amount-pop');
			void amountWrap.offsetWidth; // force reflow so animation restarts
			amountWrap.classList.add('amount-pop');
		}
	});
});
