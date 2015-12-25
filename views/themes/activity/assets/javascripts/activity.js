QorActivity = {
  init : function() {
    if($("#qor-activity").get(0)) {
      if($("#qor-activity").parents(".qor-slideout").get(0)) {
        this.$scoped = $("#qor-activity").parents(".qor-slideout");
      } else {
        this.$scoped = $("body");
      }
      this.initStatus();
      this.bindingEvents();
    }
  },

  initStatus : function() {
    this.appendTabsToFormContainer();
  },

  bindingEvents : function() {
    this.$scoped.on("click", ".qor-js-activity-item .qor-js-activity-edit", this.makeNoteEditable);
    $(".qor-js-activity-input input").blur(this.disableNote);
  },

  appendTabsToFormContainer : function() {
    var $formContainer = this.$scoped.find(".qor-form-container");
    var $scoped = this.$scoped.find(".qor-page__body");
    var $tabsWrap;
    if(!$formContainer.parent(".mdl-tabs").get(0)) {
      $scoped.append($(".qor-tabs-template").html());
      $tabsWrap = $scoped.find(".mdl-tabs");
      $('<div class="mdl-tabs__panel is-active" id="form-panel">').append($formContainer).appendTo($tabsWrap);
      $tabsWrap.find(".mdl-tabs__tab-bar").append('<a href="#tab-form-panel" class="qor-js-action-tab mdl-tabs__tab is-active">' + $(".mdl-layout-title").text() + '</a>');
    } else {
      $tabsWrap = $scoped.find(".mdl-tabs");
    }
    $tabsWrap.find(".mdl-tabs__tab-bar").append($(".qor-tabs-tab-template").html());
    $tabsWrap.append($(".qor-tabs-panel-template").html());
    $(".qor-tabs-template").remove();
    $(".qor-tabs-tab-template").remove();
    $(".qor-tabs-panel-template").remove();
  },

  makeNoteEditable : function() {
    var $item = $(this).parents(".qor-js-activity-item");
    $item.addClass("is-edit");
    $item.find(".qor-js-activity-note").hide();
    $item.find(".qor-js-activity-input").show();
    $(this).hide();
    $item.find(".qor-js-activity-save").show();
    $item.find(".qor-js-activity-input input").focus();
  },

  disableNote : function() {
    setTimeout(function() {
      var $editingItem = $(".qor-js-activity-item.is-edit");
      $editingItem.remove("is-edit");
      $editingItem.find(".qor-js-activity-note").show();
      $editingItem.find(".qor-js-activity-input").hide();
      $editingItem.find(".qor-js-activity-edit").show();
      $editingItem.find(".qor-js-activity-save").hide();
    }, 200);
  }
}
