function SWANSalt(element) {
            
    const gridId = 'salt-grid';
    const imgIdPrefix = 'salt-';
    const colIdPrefix = 'salt-col-';

    var selected = [];
    var indicators = [];

    function value() {
        if(selected.length == 4){
            var b1 = (selected[0] << 4) | selected[1];
            var b2 = (selected[2] << 4) | selected[3];

            return new Uint8Array([b1, b2]);
        }
    }

    function complete() {
        var grid = document.getElementById(gridId);
        grid.classList.add('complete');
    }

    function reset() {
        var grid = document.getElementById(gridId);
        grid.classList.remove('complete');
        selected = [];
        indicators.forEach(i => {
            i.parentNode.removeChild(i);
        });
        indicators = [];
    }

    function getIndicatorClass(value){
        switch (value) {
            case 1:
                return "top-left";
            case 2:
                return "top-right";
            case 3:
                return "bottom-left";
            case 4:
                return "bottom-right";
            default:
                throw `Invalid value: ${value}`;
        }
    }

    function updateIndicator(value) {
        var element = document.getElementById(colIdPrefix + value);
        var item = selected.length;
        
        var indicator = document.createElement('div');
        indicator.classList.add(getIndicatorClass(item));
        indicator.innerHTML = item;
        
        indicators.push(indicator);
        element.appendChild(indicator);

    }

    function add(event) {
        var value = parseInt(event.target.getAttribute("data-value"));
        
        if (selected.length < 4) {
            selected.push(value);
            updateIndicator(value);
        }
        if (selected.length == 4){
            complete();
        }
    }

    async function salt() {
        var animals = await fetch('/animals.json').then(res => res.json());

        var txt = `<div id="${gridId}" class="row">`;
        for (var i = 0; i < animals.length && i < 16; i++) {
            var animal = animals[i];

            txt += `<div id="${colIdPrefix + i}" class="col"><img id="${imgIdPrefix + i}" src="/${animal.Filename}" data-value="${i}"></div>`;

            if ((i + 1) % 4 == 0) {
                txt += '<div class="w-100"></div>';
            }
        }
        txt += '</div>'

        element.innerHTML = txt;

        for (var i = 0; i < 16; i++) {
            var animal = document.getElementById(imgIdPrefix + i);

            animal.addEventListener('click', add);
        }
    }

    //#region public methods and getters

    this.reset = function () {
        reset();
    }

    Object.defineProperty(this, 'value', { 
        get: function() { return value(); } 
    });

    //#endregion

    // generate grid;
    salt();
}