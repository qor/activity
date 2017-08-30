(function(factory) {
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
})(function($) {
    'use strict';

    let Mustache = window.Mustache,
        NAMESPACE = 'qor.activity',
        EVENT_ENABLE = 'enable.' + NAMESPACE,
        EVENT_DISABLE = 'disable.' + NAMESPACE,
        EVENT_CLICK = 'click.' + NAMESPACE,
        EVENT_SUBMIT = 'submit.' + NAMESPACE,
        CLASS_EDIT_NOTE = '.qor-activity__edit-button',
        CLASS_TAB_ACTIVITY = '.qor-tab__activity',
        CLASS_EDIT_NOTE_FORM = '.qor-activity__edit-note_form',
        CLASS_NEW_NOTE_FORM = '.qor-activity__new-note_form',
        CLASS_LISTS = '.qor-activity__lists';

    function QorActivity(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, QorActivity.DEFAULTS, $.isPlainObject(options) && options);
        this.init();
    }

    QorActivity.prototype = {
        constructor: QorActivity,

        init: function() {
            this.bind();
            this.initTabs();
        },

        bind: function() {
            this.$element.on(EVENT_CLICK, $.proxy(this.click, this)).on(EVENT_SUBMIT, 'form', $.proxy(this.submit, this));
            $(document).on(EVENT_CLICK, CLASS_TAB_ACTIVITY, $.proxy(this.tabClick, this));
        },

        unbind: function() {
            this.$element.off(EVENT_CLICK).off(EVENT_SUBMIT, 'form');
            $(document).off(EVENT_CLICK, CLASS_TAB_ACTIVITY);
        },

        submit: function(e) {
            var $form = $(e.target);
            var FormDatas;
            var self = this;
            var NoteTitle = $('#scroll-tab-activity').data().noteTitle;

            e.preventDefault();

            FormDatas = $form.serialize();
            $.ajax($form.prop('action'), {
                method: $form.prop('method'),
                data: FormDatas,
                dataType: 'json',
                headers: {
                    Accept: 'application/json; charset=utf-8'
                }
            }).done(function(data) {
                if (data.errors) {
                    return;
                }
                data.NoteTitle = NoteTitle;

                if ($form.is(CLASS_EDIT_NOTE_FORM)) {
                    self.hideEditForm($form);
                    $form.find('.qor-activity__list-note').html(data.Note);
                }

                if ($form.is(CLASS_NEW_NOTE_FORM)) {
                    $(CLASS_LISTS).prepend(self.renderActivityList(data));
                    self.clearForm();
                }
            });
            return false;
        },

        renderActivityList: function(data) {
            return Mustache.render(QorActivity.ACTIVITY_LIST_TEMPLATE, data);
        },

        clearForm: function() {
            $('textarea[data-toggle="qor.redactor"]').redactor('code.set', '');
            $(CLASS_NEW_NOTE_FORM).find('[name="QorResource.Content"],[name="QorResource.Note"]').val('');
        },

        click: function(e) {
            var $target = $(e.target);
            e.stopPropagation();

            if ($target.is(CLASS_EDIT_NOTE)) {
                var parents = $target.closest('.qor-activity__list');
                this.showEditForm(parents);
            }
        },

        tabClick: function() {
            var self = this;
            var activityList = $(CLASS_LISTS).find('.qor-activity__list').length;
            var NoteTitle = $('#scroll-tab-activity').data().noteTitle;

            if (activityList) {
                return;
            }

            var url = $(CLASS_TAB_ACTIVITY).data('resource-url');

            if (!activityList) {
                $.ajax({
                    url: url,
                    method: 'GET',
                    dataType: 'json',
                    headers: {
                        Accept: 'application/json; charset=utf-8'
                    },
                    success: function(data) {
                        if (data.length) {
                            $(CLASS_LISTS).html('');
                            for (var i = data.length - 1; i >= 0; i--) {
                                data[i].NoteTitle = NoteTitle;
                                $(CLASS_LISTS).append(self.renderActivityList(data[i]));
                            }
                        }
                        $(CLASS_LISTS).find('.mdl-spinner').remove();
                    }
                });
            } else {
                $(CLASS_LISTS).find('.mdl-spinner').remove();
            }
        },

        showEditForm: function(ele) {
            ele.find('.qor-activity__list-note,.qor-activity__edit-button').removeClass('show').addClass('hide');
            ele.find('.qor-activity__edit-feilds,.qor-activity__edit-save-button').removeClass('hide').addClass('show');
        },

        hideEditForm: function(ele) {
            ele.find('.qor-activity__list-note,.qor-activity__edit-button').removeClass('hide').addClass('show');
            ele.find('.qor-activity__edit-feilds,.qor-activity__edit-save-button').removeClass('show').addClass('hide');
        },

        initTabs: function() {
            if (!$('.qor-slideout.is-shown').get(0)) {
                $('.qor-page__body').append(QorActivity.CONTENT_HTML);
                $('.qor-form-container').appendTo($('#scroll-tab-form'));
                $('#scroll-tab-activity').appendTo('.mdl-layout__content');
                $('.qor-page__header .qor-tab-bar--activity-header').prependTo('.mdl-layout.qor-sliderout__activity-container');
                $('.qor-page > .qor-page__header').hide();
                $('.qor-page > .qor-page__header .qor-action-forms').prependTo('#scroll-tab-form');
                $('.qor-layout .mdl-layout__content.has-header').removeClass('has-header');
                $('#scroll-tab-activity').wrapInner('<div class="qor-form-container"></div>');
            }
        },

        destroy: function() {
            this.unbind();
            this.$element.removeData(NAMESPACE);
        }
    };

    QorActivity.CONTENT_HTML = `<div class="mdl-layout mdl-js-layout qor-sliderout__activity-container">
            <main class="mdl-layout__content qor-slideout--activity-content">
                <div class="mdl-layout__tab-panel is-active" id="scroll-tab-form"></div>
            </main>
        </div>`;

    QorActivity.DEFAULTS = {};

    QorActivity.ACTIVITY_LIST_TEMPLATE = `<div class="qor-activity__list">
            <form class="qor-activity__edit-note_form" action=[[ URL ]] method="POST">
                <input type="hidden" name="QorResource.ID" value=[[ ID ]]>
                <div class="qor-activity__list-title">
                    <strong>[[ CreatorName ]]</strong> <span>[[ Action ]]</span>
                </div>
                <div class="qor-activity__list-date">
                    [[ UpdatedAt ]]
                </div>
                <div class="qor-activity__list-content">[[ &Content ]]</div>
                <div>
                    <strong>[[ NoteTitle ]]</strong>
                    <span class="qor-activity__list-note">[[ Note ]]</span>
                    <a class="mdl-button mdl-js-button mdl-button--icon qor-activity__edit-button" href="#">
                        <i class="material-icons md-18 qor-activity__edit-button">edit</i>
                    </a>
                    <div class="mdl-textfield mdl-js-textfield qor-activity__edit-feilds">
                        <label class="mdl-textfield__label">[[ Note ]]</label>
                        <input class="mdl-textfield__input" type="text" name="QorResource.Note" value="[[ Note ]]">
                    </div>
                    <button class="mdl-button mdl-js-button mdl-button--icon qor-activity__edit-save-button" type="submit"><i class="material-icons md-24">done</i></button>
                </div>
            </form>
        </div>`;

    QorActivity.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {
                if (/destroy/.test(options)) {
                    return;
                }

                $this.data(NAMESPACE, (data = new QorActivity(this, options)));
            }

            if (typeof options === 'string' && $.isFunction((fn = data[options]))) {
                fn.apply(data);
            }
        });
    };

    // init activity html after sliderout loaded.
    $.fn.qorSliderAfterShow.qorActivityinit = function() {
        var $target = $('.qor-slideout > .qor-slideout__body');
        var $tab = $('.qor-slideout .qor-tab-bar--activity-header');
        $target.wrapInner(QorActivity.CONTENT_HTML);
        $('.qor-sliderout__activity-container').prepend($tab);
        $('.qor-slideout--activity-content').append($('.qor-slideout #scroll-tab-activity'));
    };

    $(function() {
        var selector = '[data-toggle="qor.activity"]';

        $(document)
            .on(EVENT_DISABLE, function(e) {
                QorActivity.plugin.call($(selector, e.target), 'destroy');
            })
            .on(EVENT_ENABLE, function(e) {
                QorActivity.plugin.call($(selector, e.target));
            })
            .triggerHandler(EVENT_ENABLE);
    });

    return QorActivity;
});
