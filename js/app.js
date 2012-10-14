var TaggedSocket = (function () {

  function log() {
    console.log('socket:', arguments)
  }

  var ws

  var TaggedSocket = function (host) {
    this.tags = []
    ws = new WebSocket(host)
    ws.onopen = $.proxy(this.onopen, this)
    ws.onclose = $.proxy(this.onclose, this)
    ws.onmessage = $.proxy(this.dispatch, this)
  }

  TaggedSocket.prototype = {

    onopen: log,
    onclose: log,
    onmessage: log,
    ontag: log,
    dispatch: function (evt) {
      var e = JSON.parse(evt.data)
      
      if(e.Tags) {
        e.Tags.forEach(function (tag) {
          this.dispatchTag(tag)
        }.bind(this))
      }
      this.onmessage(e)
    },

    dispatchTag: function (tag) {
      if(this.tags.indexOf(tag) == -1) {
        this.tags.push(tag)
        this.ontag(tag)
      }
    }

  }

  return TaggedSocket
})()

var conn;
var checkedTags = []

function runStream () {
  
  var msg = $("#msg");
  var log = $("#log");
  
  conn = new TaggedSocket("ws://{{$}}/ws");
  
  conn.onclose = function(evt) {
    alert("closed")
  }
  
  conn.onopen = function(evt) {
      // kill me
  }
  
  conn.onmessage = function(e) {
    $("#container").append(createEventDOM(e))
  }

  conn.ontag = addTagDOM;


  function filterView () {

    function updateTags (evt) {
      console.log(evt)
      var el = $(evt.currentTarget);
      var checked = el.is(':checked'),
          tag = el.attr('name');

      if(checked) {//add
        checkedTags.push(tag)
      } else if (checkedTags.indexOf(tag) > -1){ // remove
        var tags = []
        checkedTags.forEach(function(t) {
          if(t != tag) tags.push(t)
        })
        checkedTags = tags
      }

      toggleEvents()
    }

    function toggleEvents () {
      var tagsSelector = '.t_' + checkedTags.join(',.t_')
      console.log(tagsSelector, $(".event-view").filter(tagsSelector))
      // hide
       $(".event-view").filter(':not('+ tagsSelector +')').hide()
      // show
      $(".event-view").filter(tagsSelector).show()
    }   

    $('#tags').on('click', 'input[type=checkbox]', updateTags)
  }

  function addTagDOM (tag) {
    checkedTags.push(tag)

    var label = $('<label>').appendTo('#tags')
    var input = $('<input type="checkbox">').appendTo(label)
    label.append(tag)

    // input.check()
    input.attr('name', tag)
    input.attr('checked', true)

  }

  filterView()

  function createEventDOM (e) {
    var d = $("<div>").addClass('event-view')
    d.text(e.Desc + ' - ' + e.Tags.join(', '))
    d.addClass(e.Tags.join(' t_'))
    // d.attr('data-tags', e.Tags.join(' '))
    return d
  }
}