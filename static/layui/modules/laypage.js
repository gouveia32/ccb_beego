/**
 
 @Name : laypage 分页组件
 @License：MIT
 
 */

layui.define(function(exports){
  "use strict";
  
  var doc = document
  ,id = 'getElementById'
  ,tag = 'getElementsByTagName'
  
  //字符常量
  ,MOD_NAME = 'laypage', DISABLED = 'layui-disabled'
  
  //构造器
  ,Class = function(options){
    var that = this;
    that.config = options || {};
    that.config.index = ++laypage.index;
    that.render(true);
  };

  //判断传入的容器类型
  Class.prototype.type = function(){
    var config = this.config;
    if(typeof config.elem === 'object'){
      return config.elem.length === undefined ? 2 : 3;
    }
  };

  //分页视图
  Class.prototype.view = function(){
    var that = this
    ,config = that.config
    ,groups = config.groups = 'groups' in config ? (config.groups|0) : 5; //连续页码个数
    
    //排版
    config.layout = typeof config.layout === 'object' 
      ? config.layout 
    : ['prev', 'page', 'next'];
    
    config.count = config.count|0; //数据总数
    config.curr = (config.curr|0) || 1; //当前页

    //每页条数的选择项
    config.limits = typeof config.limits === 'object'
      ? config.limits
    : [10, 20, 30, 40, 50];
    config.limit = (config.limit|0) || 10; //默认条数
    
    //总页数
    config.pages = Math.ceil(config.count/config.limit) || 1;
    
    //当前页不能超过总页数
    if(config.curr > config.pages){
      config.curr = config.pages;
    }
    
    //连续分页个数不能低于0且不能大于总页数
    if(groups < 0){
      groups = 1;
    } else if (groups > config.pages){
      groups = config.pages;
    }
    
    config.prev = 'prev' in config ? config.prev : '&#x4E0A;&#x4E00;&#x9875;'; //texto da página anterior
    config.next = 'next' in config ? config.next : '&#x4E0B;&#x4E00;&#x9875;'; //texto da próxima página
    
    //Calcular grupo atual
    var index = config.pages > groups 
      ? Math.ceil( (config.curr + (groups > 1 ? 1 : 0)) / (groups > 0 ? groups : 1) )
    : 1
    
    //ver fragmento
    ,views = {
      //página anterior
      prev: function(){
        return config.prev 
          ? '<a href="javascript:;" class="layui-laypage-prev'+ (config.curr == 1 ? (' ' + DISABLED) : '') +'" data-page="'+ (config.curr - 1) +'">'+ config.prev +'</a>'
        : '';
      }()
      
      //número de página
      ,page: function(){
        var pager = [];
        
        //Quando a quantidade de dados é 0, o número da página não é exibido
        if(config.count < 1){
          return '';
        }
        
        //primeira página
        if(index > 1 && config.first !== false && groups !== 0){
          pager.push('<a href="javascript:;" class="layui-laypage-first" data-page="1"  title="&#x9996;&#x9875;">'+ (config.first || 1) +'</a>');
        }

        //Calcular a página inicial do grupo de números de página atual
        var halve = Math.floor((groups-1)/2) //igual número de páginas
        ,start = index > 1 ? config.curr - halve : 1
        ,end = index > 1 ? (function(){
          var max = config.curr + (groups - halve - 1);
          return max > config.pages ? config.pages : max;
        }()) : groups;
        
        //Impedir números de página consecutivos "não especificados" no último grupo
        if(end - start < groups - 1){
          start = end - groups + 1;
        }

        //separador esquerdo de saída
        if(config.first !== false && start > 2){
          pager.push('<span class="layui-laypage-spr">...</span>')
        }
        
        //saída de números de página consecutivos
        for(; start <= end; start++){
          if(start === config.curr){
            //pagina atual
            pager.push('<span class="layui-laypage-curr"><em class="layui-laypage-em" '+ (/^#/.test(config.theme) ? 'style="background-color:'+ config.theme +';"' : '') +'></em><em>'+ start +'</em></span>');
          } else {
            pager.push('<a href="javascript:;" data-page="'+ start +'">'+ start +'</a>');
          }
        }
        
        //saída saída separador direito e última página
        if(config.pages > groups && config.pages > end && config.last !== false){
          if(end + 1 < config.pages){
            pager.push('<span class="layui-laypage-spr">...</span>');
          }
          if(groups !== 0){
            pager.push('<a href="javascript:;" class="layui-laypage-last" title="&#x5C3E;&#x9875;"  data-page="'+ config.pages +'">'+ (config.last || config.pages) +'</a>');
          }
        }

        return pager.join('');
      }()
      
      //próxima página
      ,next: function(){
        return config.next 
          ? '<a href="javascript:;" class="layui-laypage-next'+ (config.curr == config.pages ? (' ' + DISABLED) : '') +'" data-page="'+ (config.curr + 1) +'">'+ config.next +'</a>'
        : '';
      }()
      
      //Dados totais
      ,count: '<span class="layui-laypage-count">Total: '+ config.count +' ítens</span>'
      
      //ítens por página
      ,limit: function(){
        var options = ['<span class="layui-laypage-limits"><select lay-ignore>'];
        layui.each(config.limits, function(index, item){
          options.push(
            '<option value="'+ item +'"'
            +(item === config.limit ? 'selected' : '') 
            +'>'+ item +' ítem/pag</option>'
          );
        });
        return options.join('') +'</select></span>';
      }()
      
      //刷新当前页
      ,refresh: ['<a href="javascript:;" data-page="'+ config.curr +'" class="layui-laypage-refresh">'
        ,'<i class="layui-icon layui-icon-refresh"></i>'
      ,'</a>'].join('')

      //Área de salto
      ,skip: function(){
        return ['<span class="layui-laypage-skip">vá para a pag'
          ,'<input type="text" min="1" value="'+ config.curr +'" class="layui-input">'
          ,'<button type="button" class="layui-laypage-btn">Ir</button>'
        ,'</span>'].join('');
      }()
    };

    return ['<div class="layui-box layui-laypage layui-laypage-'+ (config.theme ? (
      /^#/.test(config.theme) ? 'molv' : config.theme
    ) : 'default') +'" id="layui-laypage-'+ config.index +'">'
      ,function(){
        var plate = [];
        layui.each(config.layout, function(index, item){
          if(views[item]){
            plate.push(views[item])
          }
        });
        return plate.join('');
      }()
    ,'</div>'].join('');
  };

  //跳页的回调
  Class.prototype.jump = function(elem, isskip){
    if(!elem) return;
    var that = this
    ,config = that.config
    ,childs = elem.children
    ,btn = elem[tag]('button')[0]
    ,input = elem[tag]('input')[0]
    ,select = elem[tag]('select')[0]
    ,skip = function(){
      var curr = input.value.replace(/\s|\D/g, '')|0;
      if(curr){
        config.curr = curr;
        that.render();
      }
    };
    
    if(isskip) return skip();
    
    //页码
    for(var i = 0, len = childs.length; i < len; i++){
      if(childs[i].nodeName.toLowerCase() === 'a'){
        laypage.on(childs[i], 'click', function(){
          var curr = this.getAttribute('data-page')|0;
          if(curr < 1 || curr > config.pages) return;
          config.curr = curr;
          that.render();
        });
      }
    }
    
    //条数
    if(select){
      laypage.on(select, 'change', function(){
        var value = this.value;
        if(config.curr*value > config.count){
          config.curr = Math.ceil(config.count/value);
        }
        config.limit = value;
        that.render();
      });
    }
    
    //Sim
    if(btn){
      laypage.on(btn, 'click', function(){
        skip();
      });
    }
  };
  
  //输入页数字控制
  Class.prototype.skip = function(elem){
    if(!elem) return;
    var that = this, input = elem[tag]('input')[0];
    if(!input) return;
    laypage.on(input, 'keyup', function(e){
      var value = this.value
      ,keyCode = e.keyCode;
      if(/^(37|38|39|40)$/.test(keyCode)) return;
      if(/\D/.test(value)){
        this.value = value.replace(/\D/, '');
      }
      if(keyCode === 13){
        that.jump(elem, true)
      }
    });
  };

  //渲染分页
  Class.prototype.render = function(load){
    var that = this
    ,config = that.config
    ,type = that.type()
    ,view = that.view();
    
    if(type === 2){
      config.elem && (config.elem.innerHTML = view);
    } else if(type === 3){
      config.elem.html(view);
    } else {
      if(doc[id](config.elem)){
        doc[id](config.elem).innerHTML = view;
      }
    }

    config.jump && config.jump(config, load);
    
    var elem = doc[id]('layui-laypage-' + config.index);
    that.jump(elem);
    
    if(config.hash && !load){
      location.hash = '!'+ config.hash +'='+ config.curr;
    }
    
    that.skip(elem);
  };
  
  //外部接口
  var laypage = {
    //分页渲染
    render: function(options){
      var o = new Class(options);
      return o.index;
    }
    ,index: layui.laypage ? (layui.laypage.index + 10000) : 0
    ,on: function(elem, even, fn){
      elem.attachEvent ? elem.attachEvent('on'+ even, function(e){ //for ie
        e.target = e.srcElement;
        fn.call(elem, e);
      }) : elem.addEventListener(even, fn, false);
      return this;
    }
  }

  exports(MOD_NAME, laypage);
});