var TaggedSocket = (function () {

  function log() {
    console.log('socket:', arguments)
  }

  var ws

  var TaggedSocket = function (host, tagsQuery) {
    tagsQuery = tagsQuery || []
    this.tags = []
    ws = new WebSocket(host + "?tags=" + tagsQuery.join(','))
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
  
  // Unfiltered connection. Full stream
  conn = new TaggedSocket("ws://{{$}}/ws");
  // Filtered connection. Only events matching these tags
  // conn = new TaggedSocket("ws://{{$}}/ws", ['tag-1,tag-4']);
  
  conn.onclose = function(evt) {
    alert("closed")
  }
  
  conn.onopen = function(evt) {
      // kill me
  }
  
  conn.onmessage = function(e) {
    var d = createEventDOM(e)
    if(!hasCheckedTags(e)) d.hide()
    $("#container").append(d)
  }

  conn.ontag = addTagDOM;
  
  function hasCheckedTags (event) {
    if(checkedTags.length == 0) return true
    var matches = 0;
    checkedTags.forEach(function (t) {
      event.Tags.forEach(function (et) {
        if(t == et) ++matches
      })
    })
    return matches == checkedTags.length
  }

  function filterView () {

    function updateTags (evt) {
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
      if(checkedTags.length == 0) {// no filters selected. Show all.
        $(".event-view").show()
        return
      }
      var tagsSelector = '.t_' + checkedTags.join('.t_')
      console.log(tagsSelector, $(".event-view").filter(tagsSelector).length)
      // hide
      $(".event-view").filter(':not('+ tagsSelector +')').hide()
      // show
      $(".event-view").filter(tagsSelector).show()
    }   

    $('#tags').on('click', 'input[type=checkbox]', updateTags)
  }

  function addTagDOM (tag) {

    var label = $('<label>').appendTo('#tags')
    var input = $('<input type="checkbox">').appendTo(label)
    label.append(tag)

    // input.check()
    input.attr('name', tag)
    // input.attr('checked', true)

  }

  filterView()

  function createEventDOM (e) {
    var d = $("<div>").addClass('event-view')
    d.text(e.Desc + ' - ' + e.Tags.join(', '))
    d.addClass('t_' + e.Tags.join(' t_'))
    // d.attr('data-tags', e.Tags.join(' '))
    return d
  }
}