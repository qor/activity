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
    this.initTab();
  },

  bindingEvents : function() {
    this.$scoped.on("click", ".qor-page__body .mdl-tabs__tab", this.switchTab);
  },

  appendTabsToFormContainer : function() {
    var $formContainer = this.$scoped.find(".qor-form-container");
    var $scoped = this.$scoped.find(".qor-page__body");
    var $tabsWrap;
    if(!$formContainer.parent(".mdl-tabs").get(0)) {
      $scoped.append($(".qor-tabs-template").html());
      $tabsWrap = $scoped.find(".mdl-tabs");
      $('<div class="mdl-tabs__panel is-active" id="form-panel">').append($formContainer).appendTo($tabsWrap);
      $tabsWrap.find(".mdl-tabs__tab-bar").append('<a href="#form-panel" class="mdl-tabs__tab is-active">' + $(".mdl-layout-title").text() + '</a>');
    } else {
      $tabsWrap = $scoped.find(".mdl-tabs");
    }
    $tabsWrap.find(".mdl-tabs__tab-bar").append($(".qor-tabs-tab-template").html());
    $tabsWrap.append($(".qor-tabs-panel-template").html());
    $(".qor-tabs-template").remove();
    $(".qor-tabs-tab-template").remove();
    $(".qor-tabs-panel-template").remove();
  },

  initTab : function() {
    if(location.href.match(/#activity/)) {
      $.proxy(this.switchTab, this.$scoped.find(".mdl-tabs__tab[href='#activity-panel']"))();
    }
  },

  switchTab : function() {
    var $scoped = QorActivity.$scoped;
    $scoped.find(".mdl-tabs__tab").removeClass("is-active");
    $scoped.find(".mdl-tabs__panel").removeClass("is-active");
    $(this).addClass("is-active");
    $scoped.find($(this).attr("href")).addClass("is-active");
    var href = $(".mdl-tabs__tab.is-active").attr("href");
    console.info(href == "#activity-panel" ? href.replace("-panel", "") : href);
    location.hash = href == "#activity-panel" ? href.replace("-panel", "") : href;
    return false;
  }
}
