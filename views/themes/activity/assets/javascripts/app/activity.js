(function (factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as anonymous module.
    define(['jquery'], factory);
  } else if (typeof exports === 'object') {
    // Node / CommonJS
    factory(require('jquery'));
  } else {
    // Browser globals.
    factory(jQuery);
  }
})(function ($) {

  'use strict';

  var FormData = window.FormData;
  var NAMESPACE = 'qor.activity';
  var EVENT_ENABLE = 'enable.' + NAMESPACE;
  var EVENT_DISABLE = 'disable.' + NAMESPACE;
  var EVENT_CLICK = 'click.' + NAMESPACE;
  var EVENT_SUBMIT = 'submit.' + NAMESPACE;
  var CLASS_EDIT_NOTE = 'qor-activity__edit-button';
  var CLASS_EDIT_NOTE_FORM = 'qor-activity__edit-note_form';
  var CLASS_NEW_NOTE_FORM = 'qor-activity__new-note_form';
  var CLASS_SAVE_NOTE = 'qor-activity__edit-save-button';

  function QorActivity(element, options) {
    this.$element = $(element);
    this.options = $.extend({}, QorActivity.DEFAULTS, $.isPlainObject(options) && options);
    this.init();

  }

  QorActivity.prototype = {
    constructor: QorActivity,

    init: function () {
      var $this = this.$element;
      this.bind();
    },

    bind: function () {
      this.$element.
      on(EVENT_CLICK, $.proxy(this.click, this)).
      on(EVENT_SUBMIT,'form', $.proxy(this.submit, this));
    },

    submit: function(e){
      var form = e.target;
      var $form = $(e.target);
      var FormDatas;
      var _this = this;

        e.preventDefault();

        FormDatas = $form.serialize();
        $.ajax($form.prop('action'), {
          method: $form.prop('method'),
          data: FormDatas,
          dataType: 'json'
        }).done(function(data){
          console.log(data)
          if($form.hasClass(CLASS_EDIT_NOTE_FORM)){
            _this.hideEditForm($form);
            $form.find('.qor-activity__list-note').html(data.Note);
          }

          if($form.hasClass(CLASS_NEW_NOTE_FORM)){

          }
        });

      return false;

    },

    click: function (e) {
      var $target = $(e.target);
      e.stopPropagation();

      if($target.hasClass(CLASS_EDIT_NOTE)){
        var parents = $target.closest(".qor-activity__list");
        this.showEditForm(parents)
      }

    },
    showEditForm: function(ele){
      ele.find('.qor-activity__list-note,.qor-activity__edit-button').removeClass('show').addClass('hide');
      ele.find('.qor-activity__edit-feilds,.qor-activity__edit-save-button').removeClass('hide').addClass('show');
    },

    hideEditForm: function(ele){
      ele.find('.qor-activity__list-note,.qor-activity__edit-button').removeClass('hide').addClass('show');
      ele.find('.qor-activity__edit-feilds,.qor-activity__edit-save-button').removeClass('show').addClass('hide');
    }
  };

  QorActivity.CONTENT_HTML = (
    '<div class="mdl-layout mdl-js-layout qor-sliderout__activity-container">' +
      '<main class="mdl-layout__content qor-slideout--activity-content">' +
        '<div class="mdl-layout__tab-panel is-active" id="scroll-tab-form">' +
    '</div></main></div>'
  );

  QorActivity.DEFAULTS = {};

  QorActivity.plugin = function (options) {
    return this.each(function () {
      var $this = $(this);
      var data = $this.data(NAMESPACE);
      var fn;

      if (!data) {

        if (/destroy/.test(options)) {
          return;
        }

        $this.data(NAMESPACE, (data = new QorActivity(this, options)));
      }

      if (typeof options === 'string' && $.isFunction(fn = data[options])) {
        fn.apply(data);
      }
    });
  };

  // init activity html after sliderout loaded.
  $.fn.qorSliderAfterShow.qorActivityinit = function (url) {
    var progressURL = url;
    $.ajax({
      url: progressURL
    }).done(function (html) {
      var $target = $('.qor-slideout > .qor-slideout__body');
      var $content = $(html);
      var $tab = $content.find('.qor-tab-bar--activity-header');

      $target.wrapInner(QorActivity.CONTENT_HTML);

      $('.qor-sliderout__activity-container').prepend($tab);
      $('.qor-slideout--activity-content').append($content.find('#scroll-tab-activity'));
    });
  };

  $(function () {
    var selector = '[data-toggle="qor.activity"]';

    $(document).
      on(EVENT_DISABLE, function (e) {
        QorActivity.plugin.call($(selector, e.target), 'destroy');
      }).
      on(EVENT_ENABLE, function (e) {
        QorActivity.plugin.call($(selector, e.target));
      }).
      triggerHandler(EVENT_ENABLE);
  });

  return QorActivity;

});
