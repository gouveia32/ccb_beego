/**
 * Bootstrap Table Chinese translation
 * Author: Zhixin Wen<wenzhixin2010@gmail.com>
 */
(function ($) {
    'use strict';

    $.fn.bootstrapTable.locales['pt-BR'] = {
        formatLoadingMessage: function () {
            return 'Carregando os dados...';
        },
        formatRecordsPerPage: function (pageNumber) {
            return 'em cada pag ' + pageNumber + '  Reg';
        },
        formatShowingRows: function (pageFrom, pageTo, totalRows) {
            return 'Não mostrar ' + pageFrom + ' Início ' + pageTo + '  Reg，Total ' + totalRows + '  Reg';
        },
        formatSearch: function () {
            return 'Procurar';
        },
        formatNoMatches: function () {
            return 'Nenhum registro encontrado';
        },
        formatPaginationSwitch: function () {
            return 'mostrar/ocultar paginação';
        },
        formatRefresh: function () {
            return 'Liberado';
        },
        formatToggle: function () {
            return 'Entrega';
        },
        formatColumns: function () {
            return 'Coluna';
        },
        formatExport: function () {
            return 'Exportar';
        },
        formatClearFilters: function () {
            return 'filtro limpo';
        }
    };

    $.extend($.fn.bootstrapTable.defaults, $.fn.bootstrapTable.locales['pt-BR']);

})(jQuery);
