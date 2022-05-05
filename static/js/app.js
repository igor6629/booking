function Promt() {
    let toast = function (c) {
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {
        const {
            msg = "",
            text = "",
            footer = "",
        } = c;

        Swal.fire({
            icon: 'success',
            title: msg,
            text: text,
            footer: footer
        })
    }

    let error = function (c) {
        const {
            msg = "",
            text = "",
            footer = "",
        } = c;

        Swal.fire({
            icon: 'error',
            title: msg,
            text: text,
            footer: footer
        })
    }

    async function custom(c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c;

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            },
        })

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (c.callback !== undefined){
                        c.callback(result);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

function CheckAvailability(id) {
    document.getElementById('check-availability-button').addEventListener("click", function () {
        let html = `
            <div class="container">
                <form id="check-availability-form" action="" method="GET" novalidate class="needs-validation">
                    <div class="row">
                        <div class="col">
                            <div class="row" id="reservation-dates-modal">
                                <div class="col">
                                    <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival" autocomplete="off">
                                </div>
                                <div class="col">
                                    <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure" autocomplete="off">
                                </div>
                            </div>
                        </div>
                    </div>
                </form>
            </div>`

        attention.custom({
            msg: html,
            title: "Choose your dates",
            willOpen: () => {
                const elem = document.getElementById('reservation-dates-modal');
                const rp = new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: true,
                    minDate: new Date(),
                })
            },
            didOpen: () => {
                document.getElementById('start').removeAttribute('disabled')
                document.getElementById('end').removeAttribute('disabled')
            },
            callback: function (result) {
                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token", "{{.CSRFToken}}");
                formData.append("room_id", id)

                fetch("/search-availability-json", {
                    method: "post",
                    body: formData,
                })
                    .then(responce => responce.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon: 'success',
                                showConfirmButton: false,
                                msg: '<p>Room is available</p>'
                                    + '<p><a href="/book-room?id='
                                    + data.room_id
                                    + '&s='
                                    + data.start_date
                                    + '&e='
                                    + data.end_date
                                    + '" class="btn btn-primary">'
                                    + 'Book Now!</a></p>',
                            })
                        } else {
                            attention.error({msg: "No available"})
                        }
                    })
            }
        });
    });
}