QorActivity = {
  init : function() {
    if($("#qor-activity").get(0)) {
      this.initStatus();
      this.bindingEvents();
    }
  },

  initStatus : function() {
    this.appendTabsToFormContainer();
    this.switchTab();
  },

  bindingEvents : function() {
    $(".qor-page__body").on("click", ".mdl-tabs__tab", this.urlAddHash);
  },

  appendTabsToFormContainer : function() {
    var $formContainer = $(".qor-form-container");
    var $scoped = $(".qor-page__body");
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
  },

  switchTab : function() {
    if(location.href.match(/#activity-panel/)) {
      $(".qor-page__body .mdl-tabs__tab").removeClass("is-active");
      $(".qor-page__body .mdl-tabs__panel").removeClass("is-active");
      $(".qor-page__body .mdl-tabs__tab[href='#activity-panel']").addClass("is-active");
      $(".qor-page__body #activity-panel").addClass("is-active");
    }
  },

  urlAddHash : function() {
    var hash = $(".qor-page__body .mdl-tabs__tab.is-active").attr("href");
    location.hash = hash;
  }
}

$(document).ready(function() {
  QorActivity.init();
});
